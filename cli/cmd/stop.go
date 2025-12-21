package cmd

import (
	"fmt"
	"log"
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
		if stopChan != nil {
			stopChan <- true // Signal the daemon to stop
			fmt.Println("Daemon stop signal sent")
		} else {
			fmt.Println("Daemon is not running")
		}

		// Stop Docker Desktop
		fmt.Println("Stopping Docker Desktop...")
		dockerCmd := exec.Command("osascript", "-e", `quit app "Docker"`)
		if err := dockerCmd.Run(); err != nil {
			log.Printf("Warning: Failed to stop Docker Desktop: %v", err)
		}
	},
}

// init registers the stopCmd with the root command
func init() {
	rootCmd.AddCommand(stopCmd)
}
