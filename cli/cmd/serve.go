package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NathanBurns3/Focus-Tracker/internal/config"
	"github.com/NathanBurns3/Focus-Tracker/internal/server"
	"github.com/NathanBurns3/Focus-Tracker/internal/tracker"
	"github.com/spf13/cobra"
)

// ServeCmd is the command that runs as a background daemon (called by launchd)
var serveCmd = &cobra.Command{
	Use: "serve",
	Short: "Run as background daemon (called by launchd)",
	Long: "run the tracker as a persistent background daemon. This command is called by launchd, not directly by users.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		// Create stop channel for shutdown
		stopChan := make(chan bool)

		// Handle OS signals for shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Start server in goroutine
		go server.StartServer(cfg)

		// Start polling in goroutine
		go tracker.StartPolling(cfg, stopChan)

		//Wait for terminal signal
		sig := <-sigChan
		log.Printf("Received signal %v, shutting down gracefully...", sig)
        stopChan <- true
        fmt.Println("Daemon stopped")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}