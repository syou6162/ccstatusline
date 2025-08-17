# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ccstatusline is a CLI tool for customizing Claude Code's statusline using YAML configuration. It processes JSON input from Claude Code via stdin, executes configured shell commands, and outputs a formatted statusline.

## Development Commands

### Build
```bash
go build -o ccstatusline
```

### Test
```bash
# Run all tests
go test -v ./...

# Run specific test file
go test -v ./template_test.go ./template.go

# Run specific test function
go test -v -run TestProcessTemplate

# Run with coverage
go test -v -cover ./...
```

### Manual Testing
```bash
# Test with sample input
echo '{"model":{"display_name":"Claude 3.5"},"cwd":"/test","session_id":"abc123"}' | ./ccstatusline -config test-config.yaml

# Test with actual Claude Code JSON structure
./ccstatusline -config ~/.config/ccstatusline/config.yaml < test-input.json
```

## Architecture

### Core Flow
1. **Input**: Claude Code sends JSON via stdin containing session data (`session_id`, `cwd`, `model`, `transcript_path`, etc.)
2. **Configuration**: YAML config loaded from (in order): CLI flag, `$XDG_CONFIG_HOME/ccstatusline/config.yaml`, `~/.config/ccstatusline/config.yaml`
3. **Processing**: Each action in the config is processed sequentially:
   - Template expansion: `{.field}` patterns are replaced with JSON values using JQ queries
   - Command execution: The resulting string is executed as a shell command with JSON data available on stdin
   - Color application: ANSI color codes applied if specified
4. **Output**: Formatted statusline sent to stdout

### Key Design Decisions

**Simplified Action Structure**: Unlike early iterations, the current design has a single `command` field that:
- Always executes as a shell command
- Templates (`{.field}`) are expanded BEFORE execution
- Commands receive the full JSON input via stdin for complex processing

**Template System**:
- `{.field}` syntax uses gojq for full JQ query support
- Template expansion happens in `expandTemplates()` before command execution
- Commands can also process stdin JSON directly (e.g., `cat | jq -r '.session_id'`)

**Configuration Format**:
```yaml
actions:
  - name: required_name     # Required: unique identifier for action
    command: string         # Shell command with optional {.field} templates
    color: color_name       # Optional ANSI color
    cache_ttl: seconds      # Optional: cache results for N seconds (0 or unset = no cache)
separator: " | "            # Default: " | "
```

### Component Responsibilities

- **main.go**: Entry point, reads stdin, loads config, orchestrates processing
- **processor.go**: Executes actions (template expansion → command execution → color → caching)
- **template.go**: Handles `{.field}` expansion using gojq, provides `expandTemplates()` and legacy `processTemplate()`
- **config.go**: YAML parsing, path resolution, validates action names are unique and required
- **colors.go**: ANSI color code mapping
- **types.go**: Shared structs (Config, Action)
- **cache.go**: File-based caching with TTL support, XDG Base Directory compliant

## Important Implementation Details

### Template Processing (template.go)
- `expandTemplates()`: Only expands `{.field}` patterns, used by processor
- `processTemplate()`: Legacy function that also handles `$(command)` syntax, kept for test compatibility
- Both use gojq with query caching for performance
- Missing fields return empty strings (not errors)

### Command Execution (processor.go)
- Commands always receive JSON input via stdin
- Failed commands produce empty output (errors logged to stderr)
- Template expansion happens BEFORE command execution via `expandTemplates()`
- Cache checked before execution if `cache_ttl > 0`
- Results cached after successful execution if `cache_ttl > 0`
- Expired cache entries cleaned on startup

### Caching System (cache.go)
- Cache directory: `$XDG_CACHE_HOME/ccstatusline/` or `~/.cache/ccstatusline/`
- Cache files: `{action_name}.json` containing result and expiration timestamp
- Only caches when `cache_ttl` is explicitly set and greater than 0
- Automatic cleanup of expired entries on startup
- Designed to handle Claude Code's 3-second statusline refresh cycle

### Claude Code JSON Fields
Common fields available for templates:
- `session_id`: Session identifier
- `cwd`: Current working directory
- `model.display_name`: Model name (e.g., "Claude 3.5 Sonnet")
- `transcript_path`: Path to transcript JSON file
- `workspace.current_dir`, `workspace.project_dir`: Workspace paths
- `hook_event_name`: Event name (e.g., "Status")
- `version`: Claude Code version

## Testing Approach

- Each component has a corresponding `*_test.go` file
- Integration tests in `main_test.go` verify the complete pipeline
- Tests use table-driven patterns for comprehensive coverage
- Mock commands (e.g., `echo`) used to avoid external dependencies

## Common Tasks

### Adding a New Color
Add to `colorMap` in `colors.go`:
```go
"purple": "\033[35m",
```

### Debugging Template Issues
Test template expansion directly:
```bash
echo '{"session_id":"abc123","cwd":"/home/user"}' | go run . -config debug.yaml
```

### Processing Complex JSON Pipelines
For complex operations like extracting from transcript files:
```yaml
- name: transcript_session
  command: "cat | jq -r '.transcript_path' | xargs -I% cat % | jq -r '.sessionId' | tail -n 1"
```

### Working with Cache
For expensive operations (API calls, heavy processing):
```yaml
- name: github_issues
  command: "gh api /user/issues | jq '.total_count'"
  cache_ttl: 300  # Cache for 5 minutes
  color: green
```
