package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// installCmd defines the "install" command for the CLI
// This command installs the launchd plist file to set up the daemon service
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the daemon service",
	Long:  "Install the launchd plist file and set up the daemon service",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, _ := os.UserHomeDir() // Get user's home directory
		launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents") // LaunchAgents directory
		plistDest := filepath.Join(launchAgentsDir, "com.focustracker.daemon.plist") // Destination plist path

		// Get absolute path to binary
		binaryPath, err := os.Executable()
		if err != nil {
			log.Fatal("Failed to get executable path:", err)
		}

		workingDir := filepath.Dir(binaryPath)
		envPath := ""
		aliasesPath := ""

		if cwd, err := os.Getwd(); err == nil {
			envPath = findEnvPath(cwd)
			if envPath != "" {
				workingDir = filepath.Dir(envPath)
				aliasesCandidate := filepath.Join(workingDir, "cli", "internal", "config", "aliases.json")
				if _, err := os.Stat(aliasesCandidate); err == nil {
					aliasesPath = aliasesCandidate
				}
			}
		}

		envVars := []string{
			"<key>PATH</key>",
			"<string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>",
		}
		if envPath != "" {
			envVars = append(envVars,
				"<key>FOCUS_TRACKER_ENV_PATH</key>",
				"<string>"+envPath+"</string>",
			)
		}
		if aliasesPath != "" {
			envVars = append(envVars,
				"<key>FOCUS_TRACKER_ALIASES_PATH</key>",
				"<string>"+aliasesPath+"</string>",
			)
		}

		// Read template
		plistTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.focustracker.daemon</string>
    
    <key>ProgramArguments</key>
    <array>
        <string>` + binaryPath + `</string>
        <string>serve</string>
    </array>
    
    <key>KeepAlive</key>
    <true/>
    
    <key>StandardOutPath</key>
    <string>/tmp/focustracker.out.log</string>
    
    <key>StandardErrorPath</key>
    <string>/tmp/focustracker.err.log</string>
    
	<key>EnvironmentVariables</key>
	<dict>
		` + strings.Join(envVars, "\n        ") + `
	</dict>
    
	<key>WorkingDirectory</key>
	<string>` + workingDir + `</string>
</dict>
</plist>`

		// Create LaunchAgents directory if it doesn't exist
		if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
			log.Fatal("Failed to create LaunchAgents directory:", err)
		}

		// Write plist file
		if err := ioutil.WriteFile(plistDest, []byte(plistTemplate), 0644); err != nil {
            log.Fatal("Failed to write plist file:", err)
        }

		// Load the service
		loadCmd := exec.Command("launchctl", "load", plistDest)
		if err := loadCmd.Run(); err != nil {
			log.Printf("Warning: Failed to load service: %v", err)
        }

		fmt.Println("✓ Daemon service installed successfully")
        fmt.Printf("  Binary: %s\n", binaryPath)
        fmt.Printf("  Plist: %s\n", plistDest)
        fmt.Println("\nUse 'focus-tracker start' to run the daemon")
		fmt.Println("Use 'focus-tracker stop' to stop the daemon")
		fmt.Println("Use 'focus-tracker report' to generate a report")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func findEnvPath(startDir string) string {
	dir := startDir
	for {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}