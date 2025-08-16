package main

import (
	"strings"
	"testing"
)

func TestProcessorProcessOutput(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		inputData map[string]interface{}
		expected  string
		contains  []string
	}{
		{
			name: "simple output",
			config: &Config{
				Actions: []Action{
					{
						Name: "test",
						Command: Command{
							Type: "output",
							Text: "Hello World",
						},
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "Hello World",
		},
		{
			name: "output with template",
			config: &Config{
				Actions: []Action{
					{
						Name: "model",
						Command: Command{
							Type: "output",
							Text: "Model: {.model}",
						},
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{
				"model": "Claude 3.5",
			},
			expected: "Model: Claude 3.5",
		},
		{
			name: "output with color",
			config: &Config{
				Actions: []Action{
					{
						Name: "colored",
						Command: Command{
							Type:  "output",
							Text:  "Status",
							Color: "green",
						},
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "\033[32mStatus\033[0m",
		},
		{
			name: "multiple outputs with separator",
			config: &Config{
				Actions: []Action{
					{
						Name: "first",
						Command: Command{
							Type: "output",
							Text: "First",
						},
					},
					{
						Name: "second",
						Command: Command{
							Type: "output",
							Text: "Second",
						},
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "First | Second",
		},
		{
			name: "command followed by output",
			config: &Config{
				Actions: []Action{
					{
						Name: "echo_cmd",
						Command: Command{
							Type:    "command",
							Command: "echo 'test-output'",
						},
					},
					{
						Name: "show_output",
						Command: Command{
							Type: "output",
							Text: "Result: {command_output}",
						},
					},
				},
				Separator: " | ",
			},
			inputData: map[string]interface{}{},
			expected:  "Result: test-output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewProcessor(tt.inputData)
			result, err := processor.Process(tt.config)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}

			if tt.expected != "" && result != tt.expected {
				t.Errorf("Process() = %q, want %q", result, tt.expected)
			}

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Process() result does not contain %q", substr)
				}
			}
		})
	}
}

func TestProcessorWithComplexTemplate(t *testing.T) {
	config := &Config{
		Actions: []Action{
			{
				Name: "complex",
				Command: Command{
					Type: "output",
					Text: "{.user.name} - {.session_id | .[0:8]}",
				},
			},
		},
		Separator: " | ",
	}

	inputData := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
		},
		"session_id": "abcdefghijklmnop",
	}

	processor := NewProcessor(inputData)
	result, err := processor.Process(config)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	expected := "Alice - abcdefgh"
	if result != expected {
		t.Errorf("Process() = %q, want %q", result, expected)
	}
}
