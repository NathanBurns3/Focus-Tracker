package cmd

import (
	"fmt"

	"github.com/NathanBurns3/Focus-Tracker/internal/daemon"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start background daemon",
	Long:  "Start the background daemon that polls for the active application every 10 seconds.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting tracker daemon...")
		daemon.StartPolling()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
