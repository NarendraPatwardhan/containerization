package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the development container image",
	Run: func(cmd *cobra.Command, args []string) {
		// Run the build command using docker
		image, _ := cmd.Flags().GetString("image")
		tag, _ := cmd.Flags().GetString("tag")
		file := cmd.Flag("dockerfile").Value.String()
		fmt.Printf("Building image %s:%s using %s\n", image, tag, file)

		user := cmd.Flag("user").Value.String()
		password := cmd.Flag("password").Value.String()
		author := cmd.Flag("author").Value.String()
		email := cmd.Flag("email").Value.String()

		// Run the docker build command
		build := exec.Command("docker", "build",
			"--build-arg", fmt.Sprintf("USERNAME=%s", user),
			"--build-arg", fmt.Sprintf("PASSWORD=%s", password),
			"--build-arg", fmt.Sprintf("AUTHOR=%s", author),
			"--build-arg", fmt.Sprintf("EMAIL=%s", email),
			"-t", fmt.Sprintf("%s:%s", image, tag), "-f", file, ".")
		build.Stdout = cmd.OutOrStdout()
		build.Stderr = cmd.OutOrStderr()
		build.Env = os.Environ()
		build.Env = append(build.Env, "DOCKER_BUILDKIT=1")
		err := build.Run()
		if err != nil {
			fmt.Printf("Error building image %s:%s\n", image, tag)
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Optional flag for target image name
	buildCmd.Flags().StringP("image", "i", "machinelearning.one/devel", "Name of the image to build")

	// Required flag for target image tag
	buildCmd.Flags().StringP("tag", "t", "", "Tag of the image to build")
	buildCmd.MarkFlagRequired("tag")

	// Optional flag for source docker file
	buildCmd.Flags().StringP("dockerfile", "f", "Dockerfile", "Dockerfile to use for building the image")

	// Optional flag for user name
	buildCmd.Flags().StringP("user", "u", "compute", "User name to use in the image")

	// Required flag for password
	buildCmd.Flags().StringP("password", "p", "", "Password to use in the image")
	buildCmd.MarkFlagRequired("password")

	// Optional flag for author name and email (for git)
	buildCmd.Flags().StringP("author", "a", "Narendra Patwardhan", "Author name and email to use in the image")
	buildCmd.Flags().StringP("email", "e", "narendra@machinelearning.one", "Author email to use in the image")
}
