package cmd

import (
	"fmt"

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
	},
}

// init registers the stopCmd with the root command
func init() {
	rootCmd.AddCommand(stopCmd)
}
