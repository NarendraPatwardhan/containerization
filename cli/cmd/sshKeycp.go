package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// sshKeycpCmd represents the sshKeygen command
var sshKeycpCmd = &cobra.Command{
	Use:   "ssh-keycp",
	Short: "Copies the ssh key pair from host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the first argument as the container name
		name := args[0]

		// Get the shared directory from the flags
		shared, _ := cmd.Flags().GetString("shared")

		// Get the home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return errors.Join(errors.New("Error getting home directory"), err)
		}

		// If shared directory is empty, default to $HOME/Workspace/Shared/$name
		if shared == "" {
			shared = fmt.Sprintf("%s/Workspace/Shared/%s", home, name)
		}

		sshSrcDir := fmt.Sprintf("%s/.ssh", home)
		sshDestDir := fmt.Sprintf("%s/.ssh", shared)

		// Make sure the destination directory exists
		mkdir := exec.Command("mkdir", "-p", sshDestDir)
		mkdir.Stdout = io.Discard
		mkdir.Stderr = os.Stderr
		mkdir.Run()

		// Copy the ssh private key
		cp := exec.Command("cp", "-r", fmt.Sprintf("%s/id_rsa", sshSrcDir), sshDestDir)
		cp.Stdout = os.Stdout
		cp.Stderr = os.Stderr
		cp.Run()

		// Copy the ssh public key
		cp = exec.Command("cp", "-r", fmt.Sprintf("%s/id_rsa.pub", sshSrcDir), sshDestDir)
		cp.Stdout = os.Stdout
		cp.Stderr = os.Stderr
		cp.Run()

		return nil
	},
}

func init() {
	hooksCmd.AddCommand(sshKeycpCmd)

	// Optional flag for shared directory
	sshKeycpCmd.Flags().StringP("shared", "s", "", "Shared directory")
}
