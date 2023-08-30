package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	env "github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the development container image",
	RunE: func(cmd *cobra.Command, _ []string) error {
		// Load the environment variables
		err := env.Load()
		if err != nil {
			return errors.Join(errors.New("Error loading environment variables"), err)
		}
		// Run the build command using docker
		image, _ := cmd.Flags().GetString("image")
		tag, _ := cmd.Flags().GetString("tag")
		file := cmd.Flag("dockerfile").Value.String()
		fmt.Printf("Building image %s:%s using %s\n", image, tag, file)

		user := cmd.Flag("user").Value.String()

		password := os.Getenv("PASSWORD")

		if user != "root" && password == "" {
			return errors.New(
				"Password is required for non-root user, please set PASSWORD environment variable",
			)
		}

		author := cmd.Flag("author").Value.String()
		email := cmd.Flag("email").Value.String()

		opt := []string{
			"build",
			"--build-arg", fmt.Sprintf("AUTHOR=%s", author),
			"--build-arg", fmt.Sprintf("EMAIL=%s", email),
		}

		if user != "root" {
			opt = append(opt, "--build-arg", fmt.Sprintf("USERNAME=%s", user),
				"--secret", "id=password,env=PASSWORD")
		}

		opt = append(opt, "-t", fmt.Sprintf("%s:%s", image, tag), "-f", file, ".")

		// Run the docker build command
		build := exec.Command("docker", opt...)
		build.Stdout = cmd.OutOrStdout()
		build.Stderr = cmd.OutOrStderr()
		build.Env = os.Environ()
		build.Env = append(build.Env, "DOCKER_BUILDKIT=1")
		err = build.Run()
		if err != nil {
			return errors.Join(
				errors.New(fmt.Sprintf("Error building image %s:%s", image, tag)),
				err,
			)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Optional flag for target image name
	buildCmd.Flags().
		StringP("image", "i", "machinelearning.one/devel", "Name of the image to build")

	// Required flag for target image tag
	buildCmd.Flags().StringP("tag", "t", "", "Tag of the image to build")
	buildCmd.MarkFlagRequired("tag")

	// Optional flag for source docker file
	buildCmd.Flags().
		StringP("dockerfile", "f", "devel/main.Dockerfile", "Dockerfile to use for building the image")

	// Optional flag for user name
	buildCmd.Flags().StringP("user", "u", "compute", "Primary user name within the image")

	// Optional flag for author name and email (for git)
	buildCmd.Flags().
		StringP("author", "a", "Narendra Patwardhan", "Author name and email to use in the image")
	buildCmd.Flags().
		StringP("email", "e", "narendra@machinelearning.one", "Author email to use in the image")
}
