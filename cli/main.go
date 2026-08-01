package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/NathanBurns3/Focus-Tracker/cmd"
	"github.com/joho/godotenv"
)

func loadEnv() {
	var candidates []string
	if envPath := os.Getenv("FOCUS_TRACKER_ENV_PATH"); envPath != "" {
		candidates = append(candidates, envPath)
	}

	candidates = append(candidates, ".env", "../.env")
	if homeDir, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates,
			filepath.Join(homeDir, ".focus-tracker.env"),
			filepath.Join(homeDir, ".config", "focus-tracker", ".env"),
		)
	}

	for _, path := range candidates {
		if err := godotenv.Load(path); err == nil {
			return
		}
	}

	log.Println("No .env file found or error loading .env file")
}

func main() {
	// Load .env file if present
	loadEnv()
	cmd.Execute()
}
