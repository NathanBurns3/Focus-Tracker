package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

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

		// Give Docker a moment to start
        fmt.Println("Waiting for Docker to initialize...")
        time.Sleep(3 * time.Second)

        fmt.Println("Starting Focus Tracker daemon...")
        
        plistPath := os.ExpandEnv("$HOME/Library/LaunchAgents/com.focustracker.daemon.plist")
        
        // Load the service (if not already loaded)
        loadCmd := exec.Command("launchctl", "load", plistPath)
        loadCmd.Run() // Ignore error if already loaded
        
        // Start the service
        startCmd := exec.Command("launchctl", "start", "com.focustracker.daemon")
        if err := startCmd.Run(); err != nil {
            log.Fatalf("Failed to start daemon: %v\nMake sure you've run 'focus-tracker install' first", err)
        }
        
        fmt.Println("✓ Daemon started and running in background")
        fmt.Println("  View logs: tail -f /tmp/focustracker.out.log")
        fmt.Println("  Stop with: focus-tracker stop")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
