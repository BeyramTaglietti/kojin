package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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

func CreateFilesMap(folderPath string, folderLevel int, ignoredFolders []string) (Folder, error) {
	folder := Folder{Name: filepath.Base(folderPath), Files: []File{}, Folders: []Folder{}, Level: folderLevel}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return folder, err
	}

	for _, entry := range entries {
		path := filepath.Join(folderPath, entry.Name())

		if entry.IsDir() {

			if slices.Contains(ignoredFolders, entry.Name()) {
				continue
			}

			subFolder, err := CreateFilesMap(path, folderLevel+1, ignoredFolders)
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

func (rootFolder *Folder) WatchTree(folderPath string, command string, ignoredFolders []string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		newRootFolder, err := CreateFilesMap(folderPath, 0, ignoredFolders)
		if err != nil {
			fmt.Println("Error while walking through the folder:", err)
			continue
		}

		if filesChanged, fileName := compareTrees(*rootFolder, newRootFolder); filesChanged {
			fmt.Println("Files have changed:", fileName)

			cmd := exec.Command("sh", "-c", command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run the command and check for errors
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error running command: %v\n", err)
			}
		}

		rootFolder = &newRootFolder
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

func (f *Folder) PrintTree(indent string) {
	fmt.Printf("%s%s/\n", indent, f.Name)
	newIndent := indent + "  "
	for _, file := range f.Files {
		fmt.Printf("%s- %s (Modified: %s)\n", newIndent, file.Name, file.ModTime)
	}
	for _, subFolder := range f.Folders {
		subFolder.PrintTree(newIndent)
	}
}
