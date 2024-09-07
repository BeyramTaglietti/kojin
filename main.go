package main

import (
	"fmt"
	"kojin/utils"

	"github.com/fatih/color"
)

func main() {

	arguments, err := utils.GetArguments()

	if err != nil {
		fmt.Println("Error while getting arguments", err)
		return
	}

	if arguments.Debug {
		color.Yellow("Debug mode enabled")
	}
	if len(arguments.IgnoreList) > 0 {
		color.Cyan("Ignoring folders: %v", arguments.IgnoreList)
	}

	rootFolder, err := utils.CreateFilesMap(arguments.FolderPath, 0, arguments.IgnoreList)
	if err != nil {
		fmt.Println("Error while walking through the folder", err)
		return
	}

	if arguments.Debug {
		fmt.Println("Initial folder structure")
		rootFolder.PrintTree("")
	}

	rootFolder.WatchTree(arguments.FolderPath, arguments.Command, arguments.IgnoreList)
}
