package utils

import (
	"errors"
	"os"
	"strings"
)

const (
	ignore string = "--ignore"
	debug  string = "--debug"
)

type OsArguments struct {
	FolderPath string
	Command    string
	IgnoreList []string
	Debug      bool
}

func GetArguments() (OsArguments, error) {
	arguments := OsArguments{}

	if len(os.Args) < 2 {
		return OsArguments{}, errors.New("no folder path provided")
	}

	if len(os.Args) < 3 {
		return OsArguments{}, errors.New("no commmand provided")
	}

	arguments.FolderPath = os.Args[1]
	arguments.Command = os.Args[2]

	if len(os.Args) > 3 {
		ignoreList, cmdFound := getSettingValue(os.Args, ignore)

		if cmdFound {
			arguments.IgnoreList = strings.Split(ignoreList, ",")
		}

		_, cmdFound = getSettingValue(os.Args, debug)

		if cmdFound {
			arguments.Debug = true
		}

	}

	return arguments, nil
}

func getSettingValue(cmdArgs []string, setting string) (string, bool) {
	for _, arg := range cmdArgs {
		if strings.HasPrefix(arg, setting) {

			if !strings.Contains(arg, "=") {
				return "", true
			}

			return strings.Split(arg, "=")[1], true
		}
	}

	return "", false
}
