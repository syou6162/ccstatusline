package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Processor handles the processing of actions
type Processor struct {
	inputData map[string]interface{}
	cache     *Cache
}

// NewProcessor creates a new processor
func NewProcessor(inputData map[string]interface{}) *Processor {
	return &Processor{
		inputData: inputData,
		cache:     NewDefaultCache(),
	}
}

// Process processes the configuration and returns the final output
func (p *Processor) Process(config *Config) (string, error) {
	// Clean expired cache entries on startup
	if err := p.cache.CleanExpired(); err != nil {
		// Log but don't fail
		fmt.Fprintf(os.Stderr, "Warning: failed to clean expired cache: %v\n", err)
	}

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

	// Check cache if TTL is set
	if action.CacheTTL > 0 {
		if cachedOutput, ok := p.cache.Get(action.Name); ok {
			// Apply color to cached output if specified
			if action.Color != "" {
				cachedOutput = applyColor(cachedOutput, action.Color)
			}
			return cachedOutput, nil
		}
	}

	if action.Command != "" {
		// First, expand any templates in the command string
		expandedCommand := expandTemplates(action.Command, p.inputData)

		// Then execute as shell command
		cmd := exec.Command("sh", "-c", expandedCommand)

		// Provide JSON input via stdin
		inputJSON, _ := json.Marshal(p.inputData)
		cmd.Stdin = bytes.NewReader(inputJSON)

		var out bytes.Buffer
		cmd.Stdout = &out

		if err := cmd.Run(); err != nil {
			// Command failed, use empty output
			output = ""
		} else {
			output = strings.TrimSpace(out.String())
		}

		// Store in cache if TTL is set and output is not empty
		if action.CacheTTL > 0 && output != "" {
			if err := p.cache.Set(action.Name, output, action.CacheTTL); err != nil {
				// Log but don't fail
				fmt.Fprintf(os.Stderr, "Warning: failed to cache result for %s: %v\n", action.Name, err)
			}
		}
	}

	// Apply color if specified
	if action.Color != "" {
		output = applyColor(output, action.Color)
	}

	return output, nil
}
