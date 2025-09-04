package daemon

import (
	"fmt"
	"time"
)


func StartPolling() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Placeholder for AppleScript call
		fmt.Println("Polling for active application...")
	}
}