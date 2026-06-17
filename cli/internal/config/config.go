package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Config holds the configuration values for the application
type Config struct {
	DbPath        	string // Path to the database file
	PollingSeconds 	int    // Polling interval in seconds
	Aliases   		map[string]string // Map of application name aliases
}

// LoadAliasesFromCandidates tries each path and returns the first parsed aliases map.
func LoadAliasesFromCandidates(paths []string) map[string]string {
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		aliases := make(map[string]string)
		if err := json.Unmarshal(data, &aliases); err != nil {
			log.Println("Error parsing alias file:", err)
			return aliases
		}
		return aliases
	}

	log.Println("No alias file found, continuing without aliases")
	return make(map[string]string)
}

func aliasCandidates() []string {
	var candidates []string
	if aliasPath := os.Getenv("FOCUS_TRACKER_ALIASES_PATH"); aliasPath != "" {
		candidates = append(candidates, aliasPath)
	}

	if envPath := os.Getenv("FOCUS_TRACKER_ENV_PATH"); envPath != "" {
		candidates = append(candidates,
			filepath.Join(filepath.Dir(envPath), "cli", "internal", "config", "aliases.json"),
		)
	}

	candidates = append(candidates, "internal/config/aliases.json")
	if configDir, err := os.UserConfigDir(); err == nil {
		candidates = append(candidates, filepath.Join(configDir, "focus-tracker", "aliases.json"))
	}
	if homeDir, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(homeDir, ".focus-tracker.aliases.json"))
	}

	return candidates
}

// Load read configuration from environment variables and returns a Config struct
// If NEXT_PUBLIC_DB_PATH is not set, it logs a warning and runs in dry mode
func Load() *Config {
	cfg := Config{
		DbPath:         os.Getenv("NEXT_PUBLIC_DB_PATH"),
		PollingSeconds: 10,
		Aliases:        LoadAliasesFromCandidates(aliasCandidates()),
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