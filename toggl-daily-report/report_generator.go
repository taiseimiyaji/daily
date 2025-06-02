package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type ProjectSummary struct {
	Name       string
	TotalHours float64
	Tasks      map[string]float64
}

func GenerateReport(entries []TimeEntry, date time.Time, projectFilter string) string {
	// Filter entries by project if specified
	if projectFilter != "" {
		filtered := []TimeEntry{}
		for _, entry := range entries {
			if strings.Contains(strings.ToLower(entry.ProjectName), strings.ToLower(projectFilter)) {
				filtered = append(filtered, entry)
			}
		}
		entries = filtered
	}

	// Group entries by project
	projectMap := make(map[string]*ProjectSummary)
	totalDuration := float64(0)

	for _, entry := range entries {
		projectName := entry.ProjectName
		if projectName == "" {
			projectName = "その他"
		}

		if _, exists := projectMap[projectName]; !exists {
			projectMap[projectName] = &ProjectSummary{
				Name:  projectName,
				Tasks: make(map[string]float64),
			}
		}

		hours := float64(entry.Duration) / 3600.0
		projectMap[projectName].TotalHours += hours
		totalDuration += hours

		taskName := entry.Description
		if taskName == "" {
			taskName = "無題のタスク"
		}
		projectMap[projectName].Tasks[taskName] += hours
	}

	// Sort projects by total hours (descending)
	var projects []*ProjectSummary
	for _, project := range projectMap {
		projects = append(projects, project)
	}
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].TotalHours > projects[j].TotalHours
	})

	// Generate report
	var report strings.Builder
	report.WriteString(fmt.Sprintf("# 日報 %s\n\n", date.Format("2006-01-02")))
	report.WriteString("## サマリー\n")
	report.WriteString(fmt.Sprintf("- 稼働時間合計: %.2fh\n\n", totalDuration))
	report.WriteString("## プロジェクト別作業時間\n\n")

	for _, project := range projects {
		report.WriteString(fmt.Sprintf("### %s (%.2fh)\n", project.Name, project.TotalHours))

		// Sort tasks by hours (descending)
		type taskEntry struct {
			name  string
			hours float64
		}
		var tasks []taskEntry
		for name, hours := range project.Tasks {
			tasks = append(tasks, taskEntry{name, hours})
		}
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].hours > tasks[j].hours
		})

		for _, task := range tasks {
			report.WriteString(fmt.Sprintf("- %s: %.2fh\n", task.name, task.hours))
		}
		report.WriteString("\n")
	}

	return report.String()
}
