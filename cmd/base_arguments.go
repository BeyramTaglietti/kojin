package cmd

import (
	"errors"
	"os"
)

type osBaseArguments struct {
	WatchedFolder string
	Command       string
}

func getBaseArguments() (osBaseArguments, error) {
	arguments := osBaseArguments{}

	if len(os.Args) < 2 {
		return osBaseArguments{}, errors.New("no folder path provided")
	}

	if len(os.Args) < 3 {
		return osBaseArguments{}, errors.New("no commmand provided")
	}

	arguments.WatchedFolder = os.Args[1]
	arguments.Command = os.Args[2]

	return arguments, nil
}
