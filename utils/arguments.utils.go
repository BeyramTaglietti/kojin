package utils

import (
	"fmt"
	"os"
)

type OsArguments struct {
	FolderPath string
	Command    string
}

func GetArguments() OsArguments {
	arguments := OsArguments{}

	if len(os.Args) < 2 {
		fmt.Println("Please provide a folder path")
		return arguments
	}

	if len(os.Args) < 3 {
		fmt.Println("Please provide a command to run after the folder has changed")
		return arguments
	}

	arguments.FolderPath = os.Args[1]
	arguments.Command = os.Args[2]

	return arguments
}
