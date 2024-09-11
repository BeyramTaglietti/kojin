package cmd

import (
	"fmt"
	"kojin/watcher"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	ignore_cmd    = "ignore"
	frequency_cmd = "frequency"
)

var rootCmd = &cobra.Command{
	Use: "Kojin",
	Short: `Kojin allows you to watch a folder and run specified commands when that folder changes.

Please use the suggested flags to interact with the application.

Example:

kojin ./src "echo 'Source folder changed!'" --ignore="node_modules, tmp" 
kojin [folder path] [command] --ignore=[comma separated list of folders to ignore] --frequency=[frequency in milliseconds]
`,
	Run: func(cmd *cobra.Command, args []string) {

		ignoredFolders, err := cmd.Flags().GetStringSlice(ignore_cmd)
		if err != nil {
			log.Fatal("Error while getting the ignore list", err)
		}
		if len(ignoredFolders) > 0 {
			color.Cyan("Ignoring folders: %v", ignoredFolders)
		}

		frequency, err := cmd.Flags().GetInt(frequency_cmd)
		if err != nil {
			log.Fatal("Error while getting the frequency", err)
		}

		baseArguments, err := getBaseArguments()
		if err != nil {
			color.Red("Error while getting the base arguments: %v\nPlease use the --help flag to see the available options", err)
			os.Exit(1)
		}

		rootFolder, err := watcher.CreateFilesMap(baseArguments.WatchedFolder, 0, ignoredFolders)
		if err != nil {
			log.Fatal("Error while walking through the folder", err)
			return
		}

		rootFolder.WatchTree(baseArguments.WatchedFolder, baseArguments.Command, watcher.WatcherArguments{
			IgnoredFolders: ignoredFolders,
			Frequency:      frequency,
		})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringSlice(ignore_cmd, []string{}, "A list of folders to ignore")
	rootCmd.Flags().Int(frequency_cmd, 1000, "Frequency in milliseconds to check for changes")
}
