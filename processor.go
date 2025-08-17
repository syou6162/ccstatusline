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
	inputData map[string]interface{}
}

// NewProcessor creates a new processor
func NewProcessor(inputData map[string]interface{}) *Processor {
	return &Processor{
		inputData: inputData,
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
	var output string

	// If command is specified, execute it
	if action.Command != "" {
		cmd := exec.Command("sh", "-c", action.Command)
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			// Command failed, use empty output
			output = ""
		} else {
			output = strings.TrimSpace(out.String())
		}
	}

	// If text template is specified, process it
	if action.Text != "" {
		// Create context with command output
		context := make(map[string]interface{})
		for k, v := range p.inputData {
			context[k] = v
		}
		context["output"] = output

		output = processTemplate(action.Text, context)
	}

	// Apply color if specified
	if action.Color != "" {
		output = applyColor(output, action.Color)
	}

	return output, nil
}
