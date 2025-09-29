package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config holds the configuration values for the application
type Config struct {
	DbPath        	string // Path to the database file
	PollingSeconds 	int    // Polling interval in seconds
	Aliases   		map[string]string // Map of application name aliases
}

// LoadAliases reads the aliases from a JSON file and returns a map
func LoadAliases(path string) map[string]string {
	aliases := make(map[string]string)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("No alias file found, continuing without aliases")
		return aliases
	}
	if err := json.Unmarshal(data, &aliases); err != nil {
		log.Println("Error parsing alias file:", err)
	}
	return aliases
}

// Load read configuration from environment variables and returns a Config struct
// If NEXT_PUBLIC_DB_PATH is not set, it logs a warning and runs in dry mode
func Load() *Config {
	cfg := Config{
		DbPath:         os.Getenv("NEXT_PUBLIC_DB_PATH"),
		PollingSeconds: 10,
		Aliases:        LoadAliases("internal/config/aliases.json"),
	}

	if cfg.DbPath == "" {
		log.Println("WARNING: NEXT_PUBLIC_DB_PATH not set, running in dry mode.")
	}

	return &cfg
}

// ResolveAlias returns the alias for a given application name, or the original name if no alias exists
func (cfg *Config) ResolveAlias(name string) string {
	if alias, ok := cfg.Aliases[name]; ok {
		return alias
	}
	return name
}