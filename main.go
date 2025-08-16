package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	// Read JSON from stdin
	inputJSON, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	var inputData map[string]interface{}
	if err := json.Unmarshal(inputJSON, &inputData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Load config
	config, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Process actions
	processor := NewProcessor(inputData)
	output, err := processor.Process(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing: %v\n", err)
		os.Exit(1)
	}

	// Output result
	fmt.Print(output)
}
