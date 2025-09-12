package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/NathanBurns3/Focus-Tracker/internal/config"
	"github.com/NathanBurns3/Focus-Tracker/internal/db"
)

// UsageEntry represents a single usage log entry received from the Chrome extension
type UsageEntry struct {
	Domain 	string `json:"domain"` 	// Domain name of the website
	Minutes float32 `json:"minutes"` // Minutes spent on the website
}

// StartServer starts the HTTP server to listen for usage data from the Chrome extension
func StartServer(cfg *config.Config) {
	http.HandleFunc("/usage", func(w http.ResponseWriter, r *http.Request) {
		UsageHandler(w, r, cfg)
	})
	fmt.Println("Listening for Chrome extension on http://localhost:8080/usage")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// UsageHandler processes incoming usage data from the Chrome extension
func UsageHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	// Handle CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
    	w.WriteHeader(http.StatusOK)
    	return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var entries []UsageEntry
	if err := json.NewDecoder(r.Body).Decode(&entries); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Process each usage entry
	for _, entry := range entries {
        domain := entry.Domain
        if domain == "unknown" || domain == "newtab" {
            continue
        }

        aliasDomain := cfg.ResolveAlias(domain)
        fmt.Printf("Received usage: %s -> %.2f minutes\n", aliasDomain, entry.Minutes)
        db.InsertAppUsage(aliasDomain, entry.Minutes, "chrome")
    }
	w.WriteHeader(http.StatusOK)
}