package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// tmuxSetupCmd represents the sshKeygen command
var tmuxSetupCmd = &cobra.Command{
	Use:   "tmux-setup",
	Short: "Setup tmux for development",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the first argument as the container name
		name := args[0]

		// Start the container
		start := exec.Command("docker", "start", name)
		// Ignore stdout but route stderr to the host
		start.Stdout = io.Discard
		start.Stderr = os.Stderr
		start.Run()

		script := []string{
			"rm -rf ~/.config/tmux",
			"rm -rf ~/.tmux/plugins/",
			"git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm --depth 1",
			"git clone https://github.com/NarendraPatwardhan/terminalconfig ~/.config/tmux",
		}

		// Create array of exec arguments by joining the script array with exec -it name /bin/bash -c
		execArgs := []string{"exec", "-it", name, "/bin/bash", "-c", strings.Join(script, " && ")}

		// Run nvim setup
		setup := exec.Command("docker", execArgs...)
		// Route stdin, stdout, and stderr to the host
		setup.Stdin = os.Stdin
		setup.Stdout = os.Stdout
		setup.Stderr = os.Stderr
		setup.Run()

		fmt.Println("To install packages upon first run, press <ctrl>+<space> I")
	},
}

func init() {
	hooksCmd.AddCommand(tmuxSetupCmd)
}
