package daemon

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/NathanBurns3/Focus-Tracker/internal/config"
	"github.com/NathanBurns3/Focus-Tracker/internal/db"
)

// GetActiveApp returns the name of the currently focused application on macOS.
// It executes an AppleScript command via osascript to retrieve the application name.
func GetActiveApp() (string, error) {
	script := `tell application "System Events" to get name of first application process whose frontmost is true`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// StartPolling starts a ticker that runs every PollingSeconds interval
// On each tick, it polls for the active application and listens to stopChan
func StartPolling(cfg *config.Config, stopChan chan bool) {
	db.Connect(cfg.DbPath)
	defer db.Close()

	ticker := time.NewTicker(time.Duration(cfg.PollingSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			app, err := GetActiveApp() // Poll for the active application
			if err != nil {
				fmt.Println("Error getting active app:", err)
				continue
			}
			fmt.Println("Active application:", app)

			minutes := float32(cfg.PollingSeconds) / 60.0 // Convert seconds to minutes

			aliasApp := cfg.ResolveAlias(app) // Resolve alias if exists
			db.InsertAppUsage(aliasApp, minutes, "desktop") // Insert/update usage in the database
		case <- stopChan:
			fmt.Println("Stopping tracker daemon...")
			return
		}
	}
}