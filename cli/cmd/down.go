/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
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
			fmt.Println("Error stopping container")
			fmt.Println(err)
			return
		}

		// Remove the container
		fmt.Printf("Removing container %s\n", name)
		_, err = exec.Command("docker", "rm", name).Output()
		if err != nil {
			fmt.Println("Error removing container")
			fmt.Println(err)
			return
		}

		// Remove the shared directory if purge is true
		if purge {
			// If the shared directory is not empty, delete it
			if _, err := os.Stat(shared); !os.IsNotExist(err) {
				fmt.Printf("Removing shared directory %s\n", shared)
				err := os.RemoveAll(shared)
				if err != nil {
					fmt.Println("Error removing shared directory")
					fmt.Println(err)
					return
				}
			}
			// If the shared directory is empty, get the home directory and delete $HOME/Workspace/Shared/$name
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory")
				fmt.Println(err)
				return
			}
			shared = fmt.Sprintf("%s/Workspace/Shared/%s", home, name)
			if _, err := os.Stat(shared); !os.IsNotExist(err) {
				fmt.Printf("Removing shared directory %s\n", shared)
				err := os.RemoveAll(shared)
				if err != nil {
					fmt.Println("Error removing shared directory")
					fmt.Println(err)
					return
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	// Optional flag for purge
	downCmd.Flags().BoolP("purge", "p", false, "Purge the shared directory")

	// Optional flag for shared directory
	downCmd.Flags().StringP("shared", "s", "", "Shared directory")
}
