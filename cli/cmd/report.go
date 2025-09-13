package cmd

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/NathanBurns3/Focus-Tracker/internal/config"
	"github.com/NathanBurns3/Focus-Tracker/internal/db"
	"github.com/spf13/cobra"
)

// Helper to format minutes as "Xh Ym" or just minutes (rounds fractional minutes)
func formatMinutes(min float32) string {
	totalMin := int(math.Round(float64(min)))
	h := totalMin / 60
	m := totalMin % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// Helper to draw ASCII bars
func asciiBar(value float32, max float32, width int) string {
	barLen := int((value / max) * float32(width))
	if barLen < 1 && value > 0 {
		barLen = 1
	}
	return strings.Repeat("▇", barLen)
}

// ANSI color codes
const (
    reset   = "\033[0m"
    blue    = "\033[34m"
    green   = "\033[32m"
    yellow  = "\033[33m"
    magenta = "\033[35m"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate usage report",
	Long:  "Generate a report of application usage from the database.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		db.Connect(cfg.DbPath)
		defer db.Close()

		reports, err := db.GetTodayUsage()
		if err != nil {
			fmt.Println("Error fetching usage:", err)
			return
		}
		if len(reports) == 0 {
			fmt.Println("No usage data for today.")
			return
		}

		now := time.Now()
		fmt.Printf("🧠 Productivity Report – %s\n\n", now.Format("January 2, 2006"))

		// Separate apps and sites
		var apps, sites []db.UsageReport
		for _, r := range reports {
			if r.Source == "desktop" {
				apps = append(apps, r)
			} else if r.Source == "chrome" && r.Name != "unknown" && r.Name != "newtab" {
				sites = append(sites, r)
			}
		}

		sort.Slice(apps, func(i, j int) bool { return apps[i].Minutes > apps[j].Minutes })
		sort.Slice(sites, func(i, j int) bool { return sites[i].Minutes > sites[j].Minutes })

		// Max values for scaling bars
		maxApp := float32(1)
		if len(apps) > 0 {
			maxApp = apps[0].Minutes
		}
		maxSite := float32(1)
		if len(sites) > 0 {
			maxSite = sites[0].Minutes
		}

		// Print apps section
		fmt.Println("App Usage (minutes):" + reset)
        for i, a := range apps {
            color := blue
            if i == 0 {
                color = yellow // Highlight top app
            } else if i == 1 {
                color = green
            } else if i == 2 {
                color = magenta
            }
            fmt.Printf("%s%-20s%s %s%-30s%s %s\n",
                color, a.Name, reset,
                color, asciiBar(a.Minutes, maxApp, 30), reset,
                formatMinutes(a.Minutes),
            )
        }
        fmt.Println()

		// Print sites section
		fmt.Println("Top Sites (minutes):" + reset)
        for i, s := range sites {
            color := blue
            if i == 0 {
                color = yellow
            } else if i == 1 {
                color = green
            } else if i == 2 {
                color = magenta
            }
            fmt.Printf("%s%-20s%s %s%-30s%s %s\n",
                color, s.Name, reset,
                color, asciiBar(s.Minutes, maxSite, 30), reset,
                formatMinutes(s.Minutes),
            )
        }
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}
