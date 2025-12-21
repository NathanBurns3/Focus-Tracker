package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/NathanBurns3/Focus-Tracker/internal/config"
	"github.com/NathanBurns3/Focus-Tracker/internal/daemon"
	"github.com/NathanBurns3/Focus-Tracker/internal/server"
	"github.com/spf13/cobra"
)

// Channel to signal stopping the daemon
var stopChan chan bool

// startCmd defines the "start" command for the CLI
// This command starts the background daemon that polls for the active application and runs the server
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start background daemon",
	Long:  "Start the background daemon that polls for the active application every 10 seconds.",
	Run: func(cmd *cobra.Command, args []string) {
		// Start Docker Desktop
		fmt.Println("Starting Docker Desktop...")
		dockerCmd := exec.Command("open", "-a", "Docker")
		if err := dockerCmd.Run(); err != nil {
			log.Printf("Warning: Failed to start Docker Desktop: %v", err)
		}

		fmt.Println("Starting tracker daemon...")
		stopChan = make(chan bool)	// Channel to signal stopping the daemon
		cfg := config.Load()		 // Load configuration
		go server.StartServer(cfg)	// Start the server in a separate goroutine
		daemon.StartPolling(cfg, stopChan)	// Start the polling daemon
	},
}

// init registers the startCmd with the root command
func init() {
	rootCmd.AddCommand(startCmd)
}
