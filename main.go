package main

import (
	"fmt"
	"kojin/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type File struct {
	Name    string
	ModTime string
}

type Folder struct {
	Name    string
	Files   []File
	Folders []Folder
	Level   int
}

func main() {

	arguments := utils.GetArguments()

	rootFolder, err := createFilesMap(arguments.FolderPath, 0)
	if err != nil {
		fmt.Println("Error while walking through the folder", err)
		return
	}

	fmt.Println("Initial folder structure")
	printTree(rootFolder, "")

	watchTree(rootFolder, arguments.FolderPath, arguments.Command)
}

func createFilesMap(folderPath string, folderLevel int) (Folder, error) {
	folder := Folder{Name: filepath.Base(folderPath), Files: []File{}, Folders: []Folder{}, Level: folderLevel}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return folder, err
	}

	for _, entry := range entries {
		path := filepath.Join(folderPath, entry.Name())

		if entry.IsDir() {
			subFolder, err := createFilesMap(path, folderLevel+1)
			if err != nil {
				return folder, err
			}
			folder.Folders = append(folder.Folders, subFolder)
		} else {
			info, err := entry.Info()
			if err != nil {
				return folder, err
			}
			modTimeStr, _, _ := strings.Cut(info.ModTime().String(), " +")
			folder.Files = append(folder.Files, File{Name: entry.Name(), ModTime: modTimeStr})
		}
	}

	return folder, nil
}

func watchTree(rootFolder Folder, folderPath string, command string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		newRootFolder, err := createFilesMap(folderPath, 0)
		if err != nil {
			fmt.Println("Error while walking through the folder:", err)
			continue
		}

		if filesChanged, fileName := compareTrees(rootFolder, newRootFolder); filesChanged {
			fmt.Println("Files have changed:", fileName)

			cmd := exec.Command("sh", "-c", command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run the command and check for errors
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error running command: %v\n", err)
			}
		}

		rootFolder = newRootFolder
	}
}

func compareTrees(oldFolder, newFolder Folder) (bool, string) {
	if len(oldFolder.Files) != len(newFolder.Files) || len(oldFolder.Folders) != len(newFolder.Folders) {
		return true, ""
	}

	fileMap := make(map[string]string)
	for _, file := range oldFolder.Files {
		fileMap[file.Name] = file.ModTime
	}

	for _, newFile := range newFolder.Files {
		if modTime, exists := fileMap[newFile.Name]; !exists || modTime != newFile.ModTime {
			return true, newFile.Name
		}
	}

	for i := range oldFolder.Folders {
		if changed, fileName := compareTrees(oldFolder.Folders[i], newFolder.Folders[i]); changed {
			return changed, fileName
		}
	}

	return false, ""
}

func printTree(folder Folder, indent string) {
	fmt.Printf("%s%s/\n", indent, folder.Name)
	newIndent := indent + "  "
	for _, file := range folder.Files {
		fmt.Printf("%s- %s (Modified: %s)\n", newIndent, file.Name, file.ModTime)
	}
	for _, subFolder := range folder.Folders {
		printTree(subFolder, newIndent)
	}
}
