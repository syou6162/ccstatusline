# ccstatusline

A CLI tool for simplifying Claude Code's statusline customization with YAML configuration, inspired by [cchook](https://github.com/syou6162/cchook).

## Overview

ccstatusline makes it easy to customize Claude Code's statusline using YAML configuration instead of complex shell scripts. It provides template syntax for flexible data access and supports command execution with colored output.

## Features

- 🎨 **YAML Configuration**: Clean, readable multi-line configuration
- 📝 **Template Syntax**: `{.field}` for JSON data access and `$(command)` for shell commands
- 🌈 **Color Support**: ANSI color codes for enhanced readability
- 📂 **XDG Compliant**: Follows XDG Base Directory specification

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
  - command: "🤖 {.model.display_name}"
    color: cyan

  # Show Git branch
  - command: "($(git branch --show-current 2>/dev/null || echo 'no-git'))"
    color: green

  # Display current directory
  - command: "📁 {.cwd | split(\"/\") | .[-1]}"
    color: blue

  # Show session ID (shortened)
  - command: "[{.session_id | .[0:8]}]"
    color: gray

separator: " | "
```

### 3. Result

Your statusline will display:

```
🤖 Claude 3.5 Sonnet | (main) | 📁 myproject | [abc12345]
```

## Configuration

### Structure

```yaml
actions:
  - name: string        # Action name (optional, for debugging)
    command: string     # Shell command (templates expanded before execution)
    color: string       # Color name (optional)

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

- Basic: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`
- Bright: `gray`, `bright_red`, `bright_green`, `bright_yellow`, `bright_blue`, `bright_magenta`, `bright_cyan`, `bright_white`

## Configuration Examples

### System Information

```yaml
actions:
  - command: "💻 $(hostname -s)"
    color: magenta

  - command: "🕐 $(date +%H:%M)"
    color: yellow

separator: " | "
```

### Development Environment

```yaml
actions:
  - command: "Node: $(node -v 2>/dev/null | cut -c2- || echo 'N/A')"
    color: green

  - command: "Python: $(python3 --version 2>/dev/null | cut -d' ' -f2 || echo 'N/A')"
    color: blue

separator: " | "
```

### Minimal Configuration

```yaml
actions:
  - command: "echo '{.cwd | split(\"/\") | .[-1]} ({.model.display_name})'"
    color: cyan
```

### Complex Command Pipeline

```yaml
actions:
  # Extract session ID from transcript file
  - command: "cat | jq -r '.transcript_path' | xargs -I% cat % | jq -r '.sessionId' | tail -n 1"
    color: yellow

  # Process multiple fields with jq
  - command: "cat | jq -r '[.model.display_name, .cwd] | join(\" in \")'"
    color: cyan
```

## Configuration File Location

The configuration file is searched in the following order:

1. Path specified with `-config` flag
2. `$XDG_CONFIG_HOME/ccstatusline/config.yaml`
3. `~/.config/ccstatusline/config.yaml` (default)

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

### Colors not displaying

- Terminal support: Ensure your terminal supports ANSI color codes
- Claude Code settings: Check that Claude Code is configured to display colors

### Command output is empty

- Shell availability: Commands are executed with `sh -c`
- Error handling: Commands that fail will result in empty output
- Use `2>/dev/null` to suppress error messages in commands

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
├── config.go        # Configuration loading
├── types.go         # Type definitions
├── template.go      # Template processing
├── processor.go     # Action processing
├── colors.go        # ANSI color codes
└── *_test.go        # Test files
```

## License

MIT

## Acknowledgments

Inspired by [cchook](https://github.com/syou6162/cchook) - a similar tool for Claude Code hooks.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
