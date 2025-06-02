package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var (
		dateStr    string
		project    string
		configPath string
		help       bool
	)

	flag.StringVar(&dateStr, "date", "", "Target date for report (YYYY-MM-DD)")
	flag.StringVar(&dateStr, "d", "", "Target date for report (YYYY-MM-DD) (shorthand)")
	flag.StringVar(&project, "project", "", "Filter by project name")
	flag.StringVar(&project, "p", "", "Filter by project name (shorthand)")
	flag.StringVar(&configPath, "config", "", "Config file path")
	flag.StringVar(&configPath, "c", "", "Config file path (shorthand)")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help (shorthand)")

	flag.Parse()

	if help {
		fmt.Println("Toggl Daily Report Generator")
		fmt.Println("\nUsage:")
		fmt.Println("  toggl-daily-report [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  --date, -d    Target date for report (YYYY-MM-DD). Default: today")
		fmt.Println("  --project, -p Filter by project name")
		fmt.Println("  --config, -c  Config file path. Default: .toggl-daily-report.json in the same directory as the executable")
		fmt.Println("  --help, -h    Show this help message")
		os.Exit(0)
	}

	targetDate := time.Now()
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid date format. Use YYYY-MM-DD\n")
			os.Exit(1)
		}
		targetDate = parsedDate
	}

	if configPath == "" {
		execPath, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot find executable path: %v\n", err)
			os.Exit(1)
		}
		execDir := filepath.Dir(execPath)
		configPath = filepath.Join(execDir, ".toggl-daily-report.json")
	}

	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	client := NewTogglClient(config.APIToken)

	entries, err := client.GetTimeEntries(targetDate, config.WorkspaceID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to fetch time entries: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No time entries found for %s\n", targetDate.Format("2006-01-02"))
		os.Exit(1)
	}

	report := GenerateReport(entries, targetDate, project)
	fmt.Print(report)
}

type Config struct {
	APIToken    string `json:"api_token"`
	WorkspaceID string `json:"workspace_id"`
	DateFormat  string `json:"date_format"`
}

func loadConfig(path string) (*Config, error) {
	config := &Config{}

	apiToken := os.Getenv("TOGGL_API_TOKEN")
	if apiToken != "" {
		config.APIToken = apiToken
	}

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	if config.APIToken == "" {
		return nil, fmt.Errorf("API token not found. Set TOGGL_API_TOKEN environment variable or configure in %s", path)
	}

	return config, nil
}
