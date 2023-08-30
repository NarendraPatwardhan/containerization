package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"machinelearning.one/flux/info"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints extended information about usage",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		// Print the embedded README.md file
		fmt.Println(string(info.Readme))
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
