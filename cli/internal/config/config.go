package config

import (
	"log"
	"os"
)

// Config holds the configuration values for the application
type Config struct {
	DbPath        	string // Path to the database file
	PollingSeconds 	int    // Polling interval in seconds
}

// Load read configuration from environment variables and returns a Config struct
// If DB_PATH is not set, it logs a warning and runs in dry mode
func Load() *Config {
	cfg := Config{
		DbPath:         os.Getenv("DB_PATH"),
		PollingSeconds: 10,
	}

	if cfg.DbPath == "" {
		log.Println("WARNING: DB_PATH not set, running in dry mode.")
	}

	return &cfg
}