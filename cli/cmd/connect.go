package cmd

import (
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a development container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the first argument as the container name
		name := args[0]

		// Start the container
		start := exec.Command("docker", "start", name)
		// Ignore stdout but route stderr to the host
		start.Stdout = io.Discard
		start.Stderr = os.Stderr
		start.Run()

		// Connect to the container
		connect := exec.Command("docker", "exec", "-it", name, "/bin/bash")
		// Route stdin, stdout, and stderr to the host
		connect.Stdin = os.Stdin
		connect.Stdout = os.Stdout
		connect.Stderr = os.Stderr
		connect.Run()
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
