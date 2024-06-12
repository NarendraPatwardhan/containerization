package cmd

import (
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// sshKeygenCmd represents the sshKeygen command
var sshKeygenCmd = &cobra.Command{
	Use:   "ssh-keygen",
	Short: "Create a new ssh key pair",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the first argument as the container name
		name := args[0]

		// Get the comment flag
		comment, _ := cmd.Flags().GetString("comment")

		// Start the container
		start := exec.Command("docker", "start", name)
		// Ignore stdout but route stderr to the host
		start.Stdout = io.Discard
		start.Stderr = os.Stderr
		start.Run()

		// Run ssh-keygen
		connect := exec.Command("docker", "exec", "-it", name, "ssh-keygen", "-t", "ed25519", "-C", comment)
		// Route stdin, stdout, and stderr to the host
		connect.Stdin = os.Stdin
		connect.Stdout = os.Stdout
		connect.Stderr = os.Stderr
		connect.Run()
	},
}

func init() {
	hooksCmd.AddCommand(sshKeygenCmd)

	// Optional flag for ssh-keygen comment
	sshKeygenCmd.Flags().StringP("comment", "c", "machinelearning.one/devel", "Key Identifier")
}
