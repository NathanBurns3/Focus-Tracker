package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// stopCmd defines the "stop" command for the CLI
// This command is intended to stop the background daemon
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop background daemon",
	Long:  "Stop the background daemon that polls for the active application every 10 seconds.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping Focus Tracker daemon...")

		// Stop the service
		if err := exec.Command("launchctl", "stop", "com.focustracker.daemon").Run(); err != nil {
            log.Printf("Warning: Failed to stop daemon: %v", err)
        }

		// Unload so it doesn't auto-restart
		plistPath := os.ExpandEnv("$HOME/Library/LaunchAgents/com.focustracker.daemon.plist")
        if err := exec.Command("launchctl", "unload", plistPath).Run(); err != nil {
            log.Printf("Warning: Failed to unload daemon: %v", err)
        }

		fmt.Println("✓ Daemon stopped")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
