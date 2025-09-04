package daemon

import (
	"fmt"
	"time"

	"github.com/NathanBurns3/Focus-Tracker/internal/config"
)

// StartPolling starts a ticker that runs every PollingSeconds interval
// On each tick, it polls for the active application
func StartPolling(cfg *config.Config) {
	ticker := time.NewTicker(time.Duration(cfg.PollingSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// TODO: Placeholder for AppleScript call to get the active application
		fmt.Println("Polling for active application...")
	}
}