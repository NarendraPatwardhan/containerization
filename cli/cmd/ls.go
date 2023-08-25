package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all flux-compatible containers",
	// Add validation that exactly one arg is present and is either "images" or "containers"
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg, received %d", len(args))
		}
		if args[0] != "images" && args[0] != "containers" {
			return fmt.Errorf("arg must be either 'images' or 'containers', received %s", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		kind := args[0]
		if kind == "containers" {
			ls([]string{"docker", "ps", "-a"}, 1, "machinelearning.one")

		}
		if kind == "images" {
			ls([]string{"docker", "images"}, 0, "machinelearning.one")
		}
	},
}

func ls(cmd []string, col int, prefix string) {
	sh := exec.Command(cmd[0], cmd[1:]...)
	output, _ := sh.Output()
	lines := strings.Split(string(output), "\n")
	fmt.Println(lines[0])
	for _, line := range lines[1:] {
		if line != "" {
			columns := strings.Fields(line)
			if strings.HasPrefix(columns[col], prefix) {
				fmt.Println(line)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
