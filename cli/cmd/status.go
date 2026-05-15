package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// statusCmd defines the "status" command for the CLI
// This command checks and displays the status of the daemon service
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check daemon service status",
	Long:  "Check and display the status of the background daemon service.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check the status of the daemon service using launchctl
		out, err := exec.Command("launchctl", "list", "com.focustracker.daemon").Output()
		if err != nil {
            fmt.Println("X Daemon is not running")
            fmt.Println("Run 'focus-tracker start' to start it")
            return
        }
		if strings.Contains(string(out), "com.focustracker.daemon") {
            fmt.Println("✓ Daemon is running")
            fmt.Println("  View logs: tail -f /tmp/focustracker.out.log")
            fmt.Println("  Stop with: focus-tracker stop")
        } else {
            fmt.Println("X Daemon is not running")
            fmt.Println("   Run 'focus-tracker start' to start it")
        }
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}