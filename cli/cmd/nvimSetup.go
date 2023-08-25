package cmd

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// nvimSetupCmd represents the sshKeygen command
var nvimSetupCmd = &cobra.Command{
	Use:   "nvim-setup",
	Short: "Setup nvim for development",
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
			"rm -rf ~/.config/nvim",
			"rm -rf ~/.local/share/nvim",
			"git clone https://github.com/NvChad/NvChad ~/.config/nvim --depth 1",
			"git clone https://github.com/NarendraPatwardhan/editorconfig ~/.config/nvim/lua/custom",
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
	},
}

func init() {
	hooksCmd.AddCommand(nvimSetupCmd)
}
