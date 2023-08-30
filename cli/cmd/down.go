package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Spin down a development container",
	// Accept exactly one argument
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the first argument as the container name
		name := args[0]

		// Get the purge flag
		purge, _ := cmd.Flags().GetBool("purge")

		// Get the shared directory from the flags
		shared, _ := cmd.Flags().GetString("shared")

		// Stop the container
		fmt.Printf("Stopping container %s\n", name)
		_, err := exec.Command("docker", "stop", name).Output()
		if err != nil {
			return errors.Join(errors.New("Error stopping container"), err)
		}

		// Remove the container
		fmt.Printf("Removing container %s\n", name)
		_, err = exec.Command("docker", "rm", name).Output()
		if err != nil {
			return errors.Join(errors.New("Error removing container"), err)
		}

		// Remove the shared directory if purge is true
		if purge {
			// If the shared directory is not empty, delete it
			if _, err := os.Stat(shared); !os.IsNotExist(err) {
				fmt.Printf("Removing shared directory %s\n", shared)
				err := os.RemoveAll(shared)
				if err != nil {
					return errors.Join(errors.New("Error removing shared directory"), err)
				}
			}
			// If the shared directory is empty, get the home directory and delete $HOME/Workspace/Shared/$name
			home, err := os.UserHomeDir()
			if err != nil {
				return errors.Join(errors.New("Error getting home directory"), err)
			}
			shared = fmt.Sprintf("%s/Workspace/Shared/%s", home, name)
			if _, err := os.Stat(shared); !os.IsNotExist(err) {
				fmt.Printf("Removing shared directory %s\n", shared)
				err := os.RemoveAll(shared)
				if err != nil {
					return errors.Join(errors.New("Error removing shared directory"), err)
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	// Optional flag for purge
	downCmd.Flags().BoolP("purge", "p", false, "Purge the shared directory")

	// Optional flag for shared directory
	downCmd.Flags().StringP("shared", "s", "", "Shared directory")
}
