package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

			// Use the host network
			"--network", "host",
		}

		// Get the no-gpu flag from the flags
		noGPU, _ := cmd.Flags().GetBool("no-gpu")

		if !noGPU {
			opt = append(opt,
				// Use all GPUs
				"--gpus", "all,\"capabilities=compute,utility,graphics,display\"",
			)
		}

		x11, _ := cmd.Flags().GetBool("x11")

		if runtime.GOOS == "linux" {
			if !x11 {
				opt = append(opt, "-e", "XDG_RUNTIME_DIR=/tmp")
				opt = append(opt, "-e", fmt.Sprintf("WAYLAND_DISPLAY=%s", os.Getenv("WAYLAND_DISPLAY")))
				opt = append(opt, "-v", fmt.Sprintf("%s/%s:/tmp/%s", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY"), os.Getenv("WAYLAND_DISPLAY")))
			} else {
				opt = append(opt, "-e", "DISPLAY")
				opt = append(opt, "-v", "/tmp/.X11-unix:/tmp/.X11-unix")
				opt = append(opt, "-v", fmt.Sprintf("%s:/root/.Xauthority", os.Getenv("XAUTHORITY")))
			}
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

		extraArgs, _ := cmd.Flags().GetStringSlice("extra-args")
		for _, extraArg := range extraArgs {
			normalizedExtraArg := strings.Split(extraArg, " ")
			for _, neArg := range normalizedExtraArg {
				opt = append(opt, strings.TrimSpace(neArg))
			}
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

	// Optional flag to disable GPU access
	upCmd.Flags().BoolP("no-gpu", "g", false, "Disable GPU access")

	// Optional flag for X11 forwarding instead of Wayland
	upCmd.Flags().BoolP("x11", "x", false, "Use X11 forwarding instead of Wayland")

	// Optional flag for additional arguments
	upCmd.Flags().
		StringSliceP("extra-args", "a", []string{}, "Additional arguments to pass to the container")
}
