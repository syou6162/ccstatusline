package main

import (
	"fmt"
	"os"
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

	if action.Command != "" {
		// Process as template (supports both text and $(command) syntax)
		output = processTemplate(action.Command, p.inputData)
	}

	// Apply color if specified
	if action.Color != "" {
		output = applyColor(output, action.Color)
	}

	return output, nil
}
