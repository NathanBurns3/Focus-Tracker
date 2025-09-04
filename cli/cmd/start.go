package cmd

import (
	"fmt"

	"github.com/NathanBurns3/Focus-Tracker/internal/config"
	"github.com/NathanBurns3/Focus-Tracker/internal/daemon"
	"github.com/spf13/cobra"
)

// startCmd defines the "start" command for the CLI
// This command starts the background daemon that polls for the active application
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start background daemon",
	Long:  "Start the background daemon that polls for the active application every 10 seconds.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting tracker daemon...")
		cfg := config.Load()		 // Load configuration
		daemon.StartPolling(cfg)	// Start the polling daemon
	},
}

// init registers the startCmd with the root command
func init() {
	rootCmd.AddCommand(startCmd)
}
