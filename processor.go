package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/syou6162/ccstatusline/builtin"
)

// Processor handles the processing of actions
type Processor struct {
	inputData       map[string]interface{}
	cache           *Cache
	builtinRegistry *builtin.Registry
}

// NewProcessor creates a new processor
func NewProcessor(inputData map[string]interface{}) *Processor {
	return &Processor{
		inputData:       inputData,
		cache:           NewDefaultCache(),
		builtinRegistry: builtin.NewRegistry(),
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
	var err error

	// Get cwd from input data
	cwd := ""
	if cwdValue, ok := p.inputData["cwd"]; ok {
		if cwdStr, ok := cwdValue.(string); ok {
			cwd = cwdStr
		}
	}

	// Check cache if TTL is set
	if action.CacheTTL > 0 {
		if cachedOutput, ok := p.cache.GetWithCwd(cwd, action.Name); ok {
			// Apply prefix if specified and output is not empty
			if action.Prefix != "" && cachedOutput != "" {
				cachedOutput = action.Prefix + cachedOutput
			}
			// Apply color to cached output if specified
			if action.Color != "" {
				cachedOutput = applyColor(cachedOutput, action.Color)
			}
			return cachedOutput, nil
		}
	}

	// Process based on action type
	switch action.Type {
	case "builtin":
		// Execute builtin function
		if action.Function == "" {
			return "", fmt.Errorf("function is required for builtin type")
		}
		fn := p.builtinRegistry.Get(action.Function)
		if fn == nil {
			return "", fmt.Errorf("unknown builtin function: %s", action.Function)
		}
		output, err = fn.Execute(p.inputData)
		if err != nil {
			return "", nil // Return empty string on error
		}

	case "command", "":
		// Execute shell command (default behavior)
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
				// Command failed, return empty string (no prefix shown)
				return "", nil
			}

			output = strings.TrimSpace(out.String())

			// If output is empty, don't show prefix
			if output == "" {
				return "", nil
			}
		}

	default:
		return "", fmt.Errorf("unknown action type: %s", action.Type)
	}

	// Store in cache if TTL is set and output is not empty
	if action.CacheTTL > 0 && output != "" {
		if err := p.cache.SetWithCwd(cwd, action.Name, output, action.CacheTTL); err != nil {
			// Log but don't fail
			fmt.Fprintf(os.Stderr, "Warning: failed to cache result for %s: %v\n", action.Name, err)
		}
	}

	// Apply prefix if specified and output is not empty
	if action.Prefix != "" && output != "" {
		output = action.Prefix + output
	}

	// Apply color if specified and output is not empty
	if action.Color != "" && output != "" {
		output = applyColor(output, action.Color)
	}

	return output, nil
}
