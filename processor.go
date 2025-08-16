package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Processor handles the processing of actions
type Processor struct {
	context *Context
}

// NewProcessor creates a new processor
func NewProcessor(inputData map[string]interface{}) *Processor {
	return &Processor{
		context: &Context{
			InputJSON:      inputData,
			CommandOutputs: make(map[string]string),
			CurrentOutput:  "",
		},
	}
}

// Process processes the configuration and returns the final output
func (p *Processor) Process(config *Config) (string, error) {
	var outputs []string

	for _, action := range config.Actions {
		output, err := p.processAction(action)
		if err != nil {
			// Continue on error, just log it
			fmt.Fprintf(os.Stderr, "Error processing action %s: %v\n", action.Name, err)
			continue
		}
		if output != "" {
			outputs = append(outputs, output)
		}
	}

	// Join outputs with separator
	return strings.Join(outputs, config.Separator), nil
}

// processAction processes a single action
func (p *Processor) processAction(action Action) (string, error) {
	switch action.Command.Type {
	case "command":
		return p.processCommand(action)
	case "output":
		return p.processOutput(action)
	default:
		return "", fmt.Errorf("unknown action type: %s", action.Command.Type)
	}
}

// processCommand executes a shell command
func (p *Processor) processCommand(action Action) (string, error) {
	cmd := exec.Command("sh", "-c", action.Command.Command)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		// Command failed, but we'll continue with empty output
		p.context.CurrentOutput = ""
		p.context.InputJSON["command_output"] = ""
		return "", nil
	}

	output := strings.TrimSpace(out.String())
	p.context.CurrentOutput = output
	p.context.InputJSON["command_output"] = output
	p.context.CommandOutputs[action.Name] = output

	// Command type doesn't produce visible output directly
	return "", nil
}

// processOutput processes an output action
func (p *Processor) processOutput(action Action) (string, error) {
	// Process template
	text := processTemplate(action.Command.Text, p.context.InputJSON)

	// Apply color if specified
	if action.Command.Color != "" {
		text = applyColor(text, action.Command.Color)
	}

	return text, nil
}
