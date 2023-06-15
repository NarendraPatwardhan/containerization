package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Spin up a development container",
	// Accept exactly one argument
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the first argument as the container name
		name := args[0]

		// Get the image name and tag from the flags
		image, _ := cmd.Flags().GetString("image")
		tag, _ := cmd.Flags().GetString("tag")

		// Get the shared directory from the flags
		shared, _ := cmd.Flags().GetString("shared")
		if shared == "" {
			// Get the home directory
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory")
				fmt.Println(err)
				return
			}
			shared = fmt.Sprintf("%s/Workspace/Shared/%s", home, name)
		}

		// Check if the shared directory exists, if not create it
		if _, err := os.Stat(shared); os.IsNotExist(err) {
			err := os.MkdirAll(shared, 0755)
			if err != nil {
				fmt.Printf("Error creating shared directory %s\n", shared)
				fmt.Println(err)
				return
			}
		}

		// Get the docker group id using stat
		dockerStat, err := os.Stat("/var/run/docker.sock")
		if err != nil {
			fmt.Println("Error getting docker group id")
			fmt.Println(err)
			return
		}
		dockerGroupID := fmt.Sprintf("%d", dockerStat.Sys().(*syscall.Stat_t).Gid)

		// Get the user name from the flags
		user, _ := cmd.Flags().GetString("user")

		// Start the container in the background using docker
		fmt.Printf("Starting container %s\n", name)
		up := exec.Command("docker", "run", "-it", "-d",
			// Mount the docker socket
			"-v", "/var/run/docker.sock:/var/run/docker.sock",
			// Mount the shared directory
			"-v", fmt.Sprintf("%s:/home/%s", shared, user),
			// Use all GPUs
			"--gpus", "all",
			// Add the user to the docker group
			"--group-add", dockerGroupID,
			// Use the host network
			"--network", "host",
			fmt.Sprintf("--name=%s", name), fmt.Sprintf("%s:%s", image, tag))

		err = up.Run()
		if err != nil {
			fmt.Printf("Error starting container %s\n", name)
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Optional flag for the image name
	upCmd.Flags().StringP("image", "i", "machinelearning.one/devel", "Name of the image to use")

	// Optional flag for the image tag
	upCmd.Flags().StringP("tag", "t", "latest", "Tag of the image to use")

	// Optional flag for the shared directory
	upCmd.Flags().StringP("shared", "s", "", "Absolute path to the shared directory")

	// Optional flag for the user name
	upCmd.Flags().StringP("user", "u", "compute", "User name to use in the container")
}
