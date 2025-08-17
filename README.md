# ccstatusline

A CLI tool for simplifying Claude Code's statusline customization with YAML configuration.

## Overview

ccstatusline makes it easy to customize Claude Code's statusline using YAML configuration instead of complex shell scripts. It provides template syntax for flexible data access and supports command execution with colored output.

## Features

- **YAML Configuration**: Clean, readable multi-line configuration
- **Template Syntax**: `{.field}` for JSON data access using JQ queries
- **Shell Commands**: Execute any shell command with JSON data available via stdin
- **Color Support**: ANSI color codes for enhanced readability
- **Caching**: TTL-based caching to reduce load from frequent updates
- **XDG Compliant**: Follows XDG Base Directory specification

## Installation

```bash
go install github.com/syou6162/ccstatusline@latest
```

Or build from source:

```bash
git clone https://github.com/syou6162/ccstatusline
cd ccstatusline
go build -o ccstatusline
```

## Quick Start

### 1. Configure Claude Code

Add to your `.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "ccstatusline",
    "padding": 0
  }
}
```

### 2. Create Configuration

Create `~/.config/ccstatusline/config.yaml`:

```yaml
actions:
  # Display model name
  - name: model
    command: "echo '{.model.display_name}'"
    color: cyan

  # Show Git branch
  - name: git_branch
    command: "git branch --show-current 2>/dev/null || echo 'no-git'"
    color: green

  # Display current directory
  - name: current_dir
    command: "echo '{.cwd | split(\"/\") | .[-1]}'"
    color: blue

  # Show session ID (first 8 chars)
  - name: session
    command: "echo '{.session_id | .[0:8]}'"
    color: gray

separator: " | "
```

### 3. Result

Your statusline will display:

```
Claude 3.5 Sonnet | main | myproject | abc12345
```

## Configuration

### Structure

```yaml
actions:
  - name: string        # Required: unique identifier for action
    command: string     # Shell command (templates expanded before execution)
    prefix: string      # Optional prefix to prepend to command output
    color: string       # Color name (optional)
    cache_ttl: integer  # Cache TTL in seconds (optional, 0 or unset = no cache)

separator: string      # Separator between segments (default: " | ")
```

### How It Works

1. **Template Expansion**: `{.field}` syntax is expanded first using JQ queries
   - `{.session_id}` → `abc123def456`
   - `{.model.display_name}` → `Claude 3.5 Sonnet`
   - `{.cwd | split("/") | .[-1]}` → `myproject`

2. **Command Execution**: The expanded string is executed as a shell command
   - Commands receive Claude Code's JSON data via stdin
   - Simple commands: `whoami`, `date +%H:%M`
   - Complex pipelines: `cat | jq -r '.transcript_path' | xargs cat | jq -r '.sessionId'`

3. **Examples**:
   - Static text: `command: "echo 'Hello World'"`
   - With template: `command: "echo 'Model: {.model.display_name}'"`
   - Direct command: `command: "git branch --show-current"`
   - Using stdin: `command: "cat | jq -r '.session_id' | cut -c1-8"`

### Available Colors

**Foreground Colors:**
- Basic: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`
- Bright: `gray`, `bright_red`, `bright_green`, `bright_yellow`, `bright_blue`, `bright_magenta`, `bright_cyan`, `bright_white`

**Background Colors:**
- Basic: `bg_black`, `bg_red`, `bg_green`, `bg_yellow`, `bg_blue`, `bg_magenta`, `bg_cyan`, `bg_white`
- Bright: `bg_gray`, `bg_bright_red`, `bg_bright_green`, `bg_bright_yellow`, `bg_bright_blue`, `bg_bright_magenta`, `bg_bright_cyan`, `bg_bright_white`

## Configuration Examples

### System Information

```yaml
actions:
  - name: hostname
    command: "hostname -s"
    color: magenta

  - name: time
    command: "date +%H:%M"
    color: yellow

separator: " | "
```

### Development Environment

```yaml
actions:
  - name: node_version
    command: "node -v 2>/dev/null | cut -c2- || echo 'N/A'"
    color: green

  - name: python_version
    command: "python3 --version 2>/dev/null | cut -d' ' -f2 || echo 'N/A'"
    color: blue

separator: " | "
```

### Minimal Configuration

```yaml
actions:
  - name: status
    command: "echo '{.cwd | split(\"/\") | .[-1]} ({.model.display_name})'"
    color: cyan
```

### Complex Command Pipeline

```yaml
actions:
  # Extract session ID from transcript file
  - name: transcript_session
    command: "cat | jq -r '.transcript_path' | xargs -I% cat % | jq -r '.sessionId' | tail -n 1"
    color: yellow

  # Process multiple fields with jq
  - name: model_in_dir
    command: "cat | jq -r '[.model.display_name, .cwd] | join(\" in \")'"
    color: cyan
```

### With Prefix

```yaml
actions:
  - name: model
    command: "echo '{.model.display_name}'"
    prefix: "Model: "
    color: cyan

  - name: session
    command: "echo '{.session_id | .[0:8]}'"
    prefix: "ID: "
    color: gray
```

Output: `Model: Claude 3.5 Sonnet | ID: abc12345`

**Note**: If a command fails or returns empty output, the prefix is not displayed.

### With Background Colors

```yaml
actions:
  # Important status with background
  - name: environment
    command: "echo 'PRODUCTION'"
    color: bg_red

  # Warning with background
  - name: branch
    command: "git branch --show-current"
    color: bg_yellow

  # Info with bright background
  - name: model
    command: "echo '{.model.display_name}'"
    color: bg_bright_blue
```

### With Caching (for expensive operations)

```yaml
actions:
  # GitHub API call - cached for 5 minutes
  - name: github_issues
    command: "gh api /user/issues | jq '.total_count'"
    cache_ttl: 300
    color: green

  # Heavy processing - cached for 1 minute
  - name: docker_status
    command: "docker ps --format '{{.Names}}' | wc -l | xargs -I{} echo '{} containers'"
    cache_ttl: 60
    color: blue

  # No cache for frequently changing data
  - name: current_time
    command: "date +%H:%M:%S"
    color: yellow
```

## Configuration File Location

The configuration file is searched in the following order:

1. Path specified with `-config` flag
2. `$XDG_CONFIG_HOME/ccstatusline/config.yaml`
3. `~/.config/ccstatusline/config.yaml` (default)

## Cache Directory

Cache files are stored in (following XDG Base Directory specification):

- `$XDG_CACHE_HOME/ccstatusline/` if XDG_CACHE_HOME is set
- `~/.cache/ccstatusline/` (default)

## Command Line Options

```bash
ccstatusline -config /path/to/custom-config.yaml
```

## Input Data from Claude Code

ccstatusline receives JSON data from Claude Code via stdin, including:

- `session_id`: Current session identifier
- `cwd`: Current working directory path
- `model`: Model information (id, display_name)
- `workspace`: Workspace details (current_dir, project_dir)
- `hook_event_name`: Event name (e.g., "Status")
- `transcript_path`: Path to transcript JSON file
- `version`: Claude Code version
- `output_style`: Output formatting style

## Testing

Create a test configuration and run:

```bash
# Create test input
echo '{
  "model": {"display_name": "Claude 3.5 Sonnet"},
  "cwd": "/home/user/project",
  "session_id": "test123456789"
}' | ccstatusline -config test-config.yaml
```

## Troubleshooting

### Statusline not updating

- Check Claude Code settings: Ensure `statusLine` is configured correctly
- Verify executable: Make sure `ccstatusline` is in your PATH
- Test configuration: Run ccstatusline manually with test input
- Check action names: Ensure all actions have unique `name` fields

### Colors not displaying

- Terminal support: Ensure your terminal supports ANSI color codes
- Claude Code settings: Check that Claude Code is configured to display colors

### Command output is empty

- Shell availability: Commands are executed with `sh -c`
- Error handling: Commands that fail will result in empty output
- Use `2>/dev/null` to suppress error messages in commands

### Cache issues

- Check permissions: Ensure `~/.cache/ccstatusline/` is writable
- Clear cache: Remove files from cache directory if needed
- Disable cache: Set `cache_ttl: 0` or omit it to disable caching for specific actions

## Development

### Building

```bash
go build -o ccstatusline
```

### Testing

```bash
go test -v ./...
```

### Project Structure

```
ccstatusline/
├── main.go          # Entry point
├── config.go        # Configuration loading and validation
├── types.go         # Type definitions
├── template.go      # Template processing
├── processor.go     # Action processing with caching
├── colors.go        # ANSI color codes
├── cache.go         # Caching implementation
└── *_test.go        # Test files
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
