# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ccstatusline is a CLI tool that simplifies Claude Code's statusline customization using YAML configuration. It's inspired by cchook and follows similar design patterns - replacing complex shell scripts with clean YAML configuration and template syntax.

## Development Commands

### Build
```bash
go build -o ccstatusline
```

### Test
```bash
# Run all tests with verbose output
go test -v ./...

# Run specific test file
go test -v ./template_test.go

# Run specific test function
go test -v -run TestProcessTemplate

# Run with coverage
go test -v -cover ./...
```

### Manual Testing
```bash
# Test with sample input
echo '{"model":{"display_name":"Claude 3.5"},"current_working_directory":"/test","session_id":"abc123"}' | ./ccstatusline -config test-config.yaml

# Test with custom config
./ccstatusline -config ~/.config/ccstatusline/custom.yaml < test-input.json
```

## Architecture

### Data Flow
1. **Input**: JSON from Claude Code via stdin → `main.go`
2. **Configuration**: YAML file loaded by `config.go` (respects XDG_CONFIG_HOME)
3. **Processing**: `processor.go` executes actions sequentially:
   - `command` actions execute shell commands, store output in context
   - `output` actions apply templates and colors, produce visible text
4. **Template Engine**: `template.go` uses gojq for JQ-style queries on JSON data
5. **Output**: Single line to stdout with ANSI color codes

### Key Design Patterns

**Action Chain Pattern**: Actions execute sequentially, with `command` actions modifying the context (adding `command_output`) that subsequent `output` actions can reference. This allows chaining commands and formatting their output.

**Template Processing**: Uses `{.field}` syntax with full JQ query support. The special `{command_output}` placeholder references the last command's output without JQ processing.

**Configuration Structure**:
```yaml
actions:
  - name: identifier
    command:
      type: "command" | "output"
      command: "shell command"  # for type: command
      text: "template string"   # for type: output
      color: "color_name"       # for type: output
separator: " | "
```

### Testing Strategy

Each component has a corresponding `*_test.go` file with comprehensive unit tests. Integration tests in `main_test.go` verify the complete flow. When modifying:

- **Template syntax**: Update `template_test.go`
- **Action processing**: Update `processor_test.go`
- **Config loading**: Update `config_test.go`
- **Color codes**: Update `colors_test.go`

## Important Implementation Notes

### Template Engine (template.go)
- Uses github.com/itchyny/gojq for JQ query execution
- Caches compiled queries for performance
- Special handling for `{command_output}` to bypass JQ processing
- Returns empty string for missing fields (not an error)

### Processor Context (processor.go)
- Maintains `Context` struct with InputJSON and CommandOutputs
- Updates `InputJSON["command_output"]` after each command execution
- Continues processing on command failures (returns empty output)
- Logs errors to stderr but doesn't stop the pipeline

### Configuration Resolution (config.go)
- Checks paths in order: CLI flag → XDG_CONFIG_HOME → ~/.config
- Sets default separator to single space if not specified
- Returns error if config file doesn't exist (no fallback to defaults)

## GitHub Actions

### Test Workflow
`.github/workflows/test.yml` runs on every push:
- Tests against Go 1.21 and 1.22
- Runs tests with race detection
- Reports test coverage

## Common Modifications

### Adding New Color
Add to `colorMap` in `colors.go`:
```go
"new_color": "\033[XXm",
```

### Adding New Template Function
Extend JQ query support in `template.go` - gojq handles most standard JQ functions automatically.

### Adding New Action Type
1. Update `Command` struct in `types.go`
2. Add case in `processAction()` in `processor.go`
3. Add corresponding test in `processor_test.go`
