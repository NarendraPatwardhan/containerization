package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Spin up a development container",
	// Accept exactly one argument
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
				return errors.Join(errors.New("Error getting home directory"), err)
			}
			shared = fmt.Sprintf("%s/Workspace/Shared/%s", home, name)
		}

		// Check if the shared directory exists, if not create it
		if _, err := os.Stat(shared); os.IsNotExist(err) {
			err := os.MkdirAll(shared, 0755)
			if err != nil {
				return errors.Join(errors.New("Error creating shared directory"), err)
			}
		}

		// Get the user name from the flags
		user, _ := cmd.Flags().GetString("user")

		opt := []string{
			"run", "-it", "-d",
			// Use all GPUs
			"--gpus", "all",
			// Use the host network
			"--network", "host",
		}

		if runtime.GOOS == "linux" {
			opt = append(opt, "-e", "DISPLAY")
			opt = append(opt, "-v", "/tmp/.X11-unix:/tmp/.X11-unix")
			opt = append(opt, "-v", fmt.Sprintf("%s:/root/.Xauthority", os.Getenv("XAUTHORITY")))
		}

		// Mount the shared directory
		if user != "root" {
			opt = append(opt, "-v", fmt.Sprintf("%s:/home/%s", shared, user))
		} else {
			opt = append(opt, "-v", fmt.Sprintf("%s:/root", shared))
		}

		// Get the recursive profile from the flags
		recursive, _ := cmd.Flags().GetString("recursive")
		switch recursive {
		case "docker":
			// Uses Docker out of Docker setup
			// Get the docker group id using stat
			dockerStat, err := os.Stat("/var/run/docker.sock")
			if err != nil {
				return errors.Join(errors.New("Error getting docker group id"), err)
			}
			dockerGroupID := fmt.Sprintf("%d", dockerStat.Sys().(*syscall.Stat_t).Gid)
			opt = append(opt,
				// Mount the docker socket
				"-v", "/var/run/docker.sock:/var/run/docker.sock",
				// Add the user to the docker group
				"--group-add", dockerGroupID,
			)
		case "containerd":
			// Uses Containerd inside Docker setup
			// Referenced from https://github.com/containerd/containerd/discussions/5522
			opt = append(opt,
				"--privileged",
				"-v", "/var/lib/containerd",
				"--tmpfs", "/run",
			)
		case "none":
			// Do nothing
		default:
			return errors.New(
				fmt.Sprintf(
					"Unknown recursive containerization %s\n, valid options are docker, containerd and none",
					recursive,
				),
			)
		}

		opt = append(opt, fmt.Sprintf("--name=%s", name), fmt.Sprintf("%s:%s", image, tag))

		// Start the container in the background using docker
		fmt.Printf("Starting container %s\n", name)
		up := exec.Command("docker", opt...)

		err := up.Run()
		if err != nil {
			return errors.Join(errors.New("Error starting container"), err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Optional flag for the image name
	upCmd.Flags().StringP("image", "i", "machinelearning.one/devel", "Name of the image to use")

	// Optional flag for the image tag
	upCmd.Flags().StringP("tag", "t", "main", "Tag of the image to use")

	// Optional flag for the shared directory
	upCmd.Flags().StringP("shared", "s", "", "Absolute path to the shared directory")

	// Optional flag for the user name
	upCmd.Flags().StringP("user", "u", "compute", "User name to use in the container")

	// Optional flag for the recursive containerization
	upCmd.Flags().
		StringP("recursive", "r", "docker", "Recursive profile - docker, containerd or none")
}
