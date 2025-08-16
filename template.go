package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/itchyny/gojq"
)

// JQ query cache for performance
var (
	jqQueryCache = make(map[string]*gojq.Query)
	jqCacheMutex sync.RWMutex
)

// executeJQQuery executes a gojq query and returns the result as a string
func executeJQQuery(queryStr string, input interface{}) (string, error) {
	// Get query from cache or create new one
	jqCacheMutex.RLock()
	query, exists := jqQueryCache[queryStr]
	jqCacheMutex.RUnlock()

	if !exists {
		// Parse query and cache it
		var err error
		query, err = gojq.Parse(queryStr)
		if err != nil {
			return "", fmt.Errorf("invalid jq query '%s': %w", queryStr, err)
		}

		jqCacheMutex.Lock()
		jqQueryCache[queryStr] = query
		jqCacheMutex.Unlock()
	}

	// Convert input to gojq-compatible type
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input to JSON: %w", err)
	}

	var gojqInput interface{}
	if err := json.Unmarshal(inputJSON, &gojqInput); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON for gojq: %w", err)
	}

	// Execute query
	iter := query.Run(gojqInput)
	var results []interface{}

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", fmt.Errorf("jq query execution error: %w", err)
		}
		results = append(results, v)
	}

	// Convert results to string
	switch len(results) {
	case 0:
		return "", nil
	case 1:
		return jqValueToString(results[0]), nil
	default:
		// Return as JSON array for multiple results
		resultJSON, err := json.Marshal(results)
		if err != nil {
			return "", fmt.Errorf("failed to marshal jq results: %w", err)
		}
		return string(resultJSON), nil
	}
}

// jqValueToString converts a gojq result value to string
func jqValueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case nil:
		return ""
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		// For numbers and objects, output as JSON
		if result, err := json.Marshal(v); err == nil {
			return string(result)
		}
		return fmt.Sprintf("%v", v)
	}
}

// processTemplate processes template strings with {.field} syntax
func processTemplate(template string, data map[string]interface{}) string {
	// Pattern: { followed by any content and ending with }
	pattern := regexp.MustCompile(`\{([^}]+)\}`)

	return pattern.ReplaceAllStringFunc(template, func(match string) string {
		content := strings.TrimSpace(match[1 : len(match)-1]) // Remove {}

		// Special case for command_output
		if content == "command_output" {
			if output, ok := data["command_output"].(string); ok {
				return output
			}
			return ""
		}

		// Process as JQ query
		result, err := executeJQQuery(content, data)
		if err != nil {
			return fmt.Sprintf("[ERROR: %s]", err.Error())
		}
		return result
	})
}
