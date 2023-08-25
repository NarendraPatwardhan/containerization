package cmd

import (
	"github.com/spf13/cobra"
)

// hooksCmd represents the hooks command
var hooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "Run the hooks for given container",
}

func init() {
	rootCmd.AddCommand(hooksCmd)
}
