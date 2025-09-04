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
		fmt.Println("Stopping tracker daemon...")
		// TODO: Implement stop functionality
	},
}

// init registers the stopCmd with the root command
func init() {
	rootCmd.AddCommand(stopCmd)
}
