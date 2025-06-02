package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TogglClient struct {
	apiToken string
	baseURL  string
	client   *http.Client
}

type TimeEntry struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	Stop        time.Time `json:"stop"`
	Duration    int64     `json:"duration"`
	ProjectID   int64     `json:"project_id"`
	ProjectName string
	TagIDs      []int64 `json:"tag_ids"`
}

type Project struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewTogglClient(apiToken string) *TogglClient {
	return &TogglClient{
		apiToken: apiToken,
		baseURL:  "https://api.track.toggl.com/api/v9",
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *TogglClient) GetTimeEntries(date time.Time, workspaceID string) ([]TimeEntry, error) {
	if workspaceID == "" {
		// Get default workspace if not specified
		var err error
		workspaceID, err = c.getDefaultWorkspaceID()
		if err != nil {
			return nil, fmt.Errorf("failed to get default workspace: %w", err)
		}
	}

	// Use the standard Toggl API v9 for time entries
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(24 * time.Hour).Add(-1 * time.Second)

	url := fmt.Sprintf("%s/me/time_entries?start_date=%s&end_date=%s",
		c.baseURL,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.apiToken, "api_token")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %s: %s", resp.Status, string(body))
	}

	var timeEntries []struct {
		ID          int64   `json:"id"`
		Description string  `json:"description"`
		Start       string  `json:"start"`
		Stop        string  `json:"stop"`
		Duration    int64   `json:"duration"`
		ProjectID   *int64  `json:"project_id"`
		WorkspaceID int64   `json:"workspace_id"`
		TagIDs      []int64 `json:"tag_ids"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&timeEntries); err != nil {
		return nil, err
	}

	// Get project names
	projectMap := make(map[int64]string)
	if len(timeEntries) > 0 {
		var err error
		projectMap, err = c.getProjects(workspaceID)
		if err != nil {
			// Log error but continue without project names
			fmt.Printf("Warning: Failed to get project names: %v\n", err)
		}
	}

	var entries []TimeEntry
	for _, te := range timeEntries {
		entry := TimeEntry{
			ID:          te.ID,
			Description: te.Description,
			Duration:    te.Duration,
		}

		if te.ProjectID != nil {
			entry.ProjectID = *te.ProjectID
			if name, ok := projectMap[*te.ProjectID]; ok {
				entry.ProjectName = name
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (c *TogglClient) getProjects(workspaceID string) (map[int64]string, error) {
	if workspaceID == "" {
		// Get default workspace if not specified
		workspaceID, _ = c.getDefaultWorkspaceID()
	}

	url := fmt.Sprintf("%s/workspaces/%s/projects", c.baseURL, workspaceID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.apiToken, "api_token")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get projects: %s", resp.Status)
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	projectMap := make(map[int64]string)
	for _, p := range projects {
		projectMap[p.ID] = p.Name
	}

	return projectMap, nil
}

func (c *TogglClient) getDefaultWorkspaceID() (string, error) {
	url := fmt.Sprintf("%s/me", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(c.apiToken, "api_token")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	var userData struct {
		DefaultWorkspaceID int64 `json:"default_workspace_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", userData.DefaultWorkspaceID), nil
}
