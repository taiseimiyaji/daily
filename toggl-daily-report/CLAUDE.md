# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Run
```bash
# Build the CLI tool
go build -o toggl-daily-report

# Run with various options
./toggl-daily-report                                    # Today's report
./toggl-daily-report --date 2025-06-01                 # Specific date
./toggl-daily-report --project "ProjectName"           # Filter by project
./toggl-daily-report --config /path/to/config.json     # Custom config path
```

### Development
```bash
# Install dependencies (none required - uses only Go standard library)
go mod tidy

# Run directly without building
go run . --date 2025-06-01
```

## Architecture

### Project Structure
- `main.go` - CLI entry point, handles argument parsing and configuration loading
- `toggl_client.go` - Toggl API v9 client implementation
- `report_generator.go` - Business logic for aggregating time entries and generating Markdown reports

### Key Design Decisions
1. **Zero external dependencies** - Uses only Go standard library for maximum portability
2. **Configuration hierarchy** - Environment variable `TOGGL_API_TOKEN` takes precedence over config file
3. **Config file location** - `.toggl-daily-report.json` is read from the same directory as the executable (not home directory)
4. **API Integration** - Uses Toggl Track API v9 with basic authentication
5. **Output format** - Generates Japanese-language Markdown reports with time rounded to 0.1 hours

### Data Flow
1. Load configuration from environment or `.toggl-daily-report.json`
2. Create Toggl client with API token
3. Fetch time entries for specified date using `/me/time_entries` endpoint
4. Fetch project names from `/workspaces/{id}/projects` endpoint
5. Aggregate entries by project and task
6. Generate Markdown report sorted by total project hours

### Important Implementation Details
- Time entries are fetched in UTC to avoid timezone parsing issues
- Project names are fetched separately and mapped to entries by project ID
- Empty project names default to "その他" (Other)
- Empty task descriptions default to "無題のタスク" (Untitled Task)
- Hours are displayed with 1 decimal place precision (e.g., 1.5h, 2.3h)

### API Endpoints Used
- `GET /api/v9/me` - Get user info and default workspace
- `GET /api/v9/me/time_entries` - Fetch time entries with date filtering
- `GET /api/v9/workspaces/{id}/projects` - Get project information