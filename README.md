# ccstatusline

A CLI tool for simplifying Claude Code's statusline customization with YAML configuration, inspired by [cchook](https://github.com/syou6162/cchook).

## Overview

ccstatusline makes it easy to customize Claude Code's statusline using YAML configuration instead of complex shell scripts. It provides template syntax for flexible data access and supports command execution with colored output.

## Features

- 🎨 **YAML Configuration**: Clean, readable multi-line configuration
- 📝 **Template Syntax**: Simple `{.field}` syntax for accessing JSON data with full jq query support
- 🔧 **Action System**: Execute shell commands and format output
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
  - name: model
    command:
      type: output
      text: "🤖 {.model.display_name}"
      color: cyan

  # Show Git branch
  - name: git_branch
    command:
      type: command
      command: "git branch --show-current 2>/dev/null || echo 'no-git'"

  - name: git_output
    command:
      type: output
      text: " ({command_output})"
      color: green

  # Display current directory
  - name: directory
    command:
      type: output
      text: " 📁 {.current_working_directory | split(\"/\") | .[-1]}"
      color: blue

  # Show session ID (shortened)
  - name: session
    command:
      type: output
      text: " [{.session_id | .[0:8]}]"
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
    command:           # Command configuration
      type: string     # "command" or "output"
      # For type: "command"
      command: string  # Shell command to execute
      # For type: "output"
      text: string     # Text to display (supports templates)
      color: string    # Color name (optional)

separator: string      # Separator between segments (default: " ")
```

### Action Types

#### `command` Action
Executes a shell command and stores the result in `{command_output}` for subsequent actions.

```yaml
- name: git_branch
  command:
    type: command
    command: "git branch --show-current"
```

#### `output` Action
Displays text with template expansion and optional color.

```yaml
- name: show_branch
  command:
    type: output
    text: "Branch: {command_output}"
    color: green
```

### Template Syntax

Access JSON fields from Claude Code using `{.field}` syntax:

- **Simple fields**: `{.session_id}`, `{.model.display_name}`
- **Nested fields**: `{.model.id}`, `{.workspace.name}`
- **JQ filters**: `{.session_id | .[0:8]}`, `{.path | split("/") | .[-1]}`
- **Command output**: `{command_output}` (from previous command action)

### Available Colors

- Basic: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`
- Bright: `gray`, `bright_red`, `bright_green`, `bright_yellow`, `bright_blue`, `bright_magenta`, `bright_cyan`, `bright_white`

## Configuration Examples

### System Information

```yaml
actions:
  - name: hostname
    command:
      type: command
      command: "hostname -s"

  - name: host_output
    command:
      type: output
      text: "💻 {command_output}"
      color: magenta

  - name: time
    command:
      type: command
      command: "date +%H:%M"

  - name: time_output
    command:
      type: output
      text: " 🕐 {command_output}"
      color: yellow

separator: " | "
```

### Development Environment

```yaml
actions:
  - name: node_version
    command:
      type: command
      command: "node -v 2>/dev/null | cut -c2- || echo 'N/A'"

  - name: node_output
    command:
      type: output
      text: "Node: {command_output}"
      color: green

  - name: python_version
    command:
      type: command
      command: "python3 --version 2>/dev/null | cut -d' ' -f2 || echo 'N/A'"

  - name: python_output
    command:
      type: output
      text: " Python: {command_output}"
      color: blue

separator: " | "
```

### Minimal Configuration

```yaml
actions:
  - name: simple
    command:
      type: output
      text: "📍 {.current_working_directory | split(\"/\") | .[-1]} ({.model.display_name})"
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
- `current_working_directory`: Current working directory path
- `model`: Model information (id, display_name)
- `workspace`: Workspace details
- `claude_code_version`: Claude Code version
- `output_style`: Output formatting style

## Testing

Create a test configuration and run:

```bash
# Create test input
echo '{
  "model": {"display_name": "Claude 3.5 Sonnet"},
  "current_working_directory": "/home/user/project",
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
