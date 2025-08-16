package main

import (
	"testing"
)

func TestProcessTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "simple field access",
			template: "Model: {.model}",
			data: map[string]interface{}{
				"model": "Claude 3.5",
			},
			expected: "Model: Claude 3.5",
		},
		{
			name:     "nested field access",
			template: "Name: {.user.name}",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
				},
			},
			expected: "Name: John",
		},
		{
			name:     "command output",
			template: "Branch: {command_output}",
			data: map[string]interface{}{
				"command_output": "main",
			},
			expected: "Branch: main",
		},
		{
			name:     "jq filter - slice",
			template: "ID: {.session_id | .[0:8]}",
			data: map[string]interface{}{
				"session_id": "abc123def456ghi789",
			},
			expected: "ID: abc123de",
		},
		{
			name:     "jq filter - split and last",
			template: "Dir: {.path | split(\"/\") | .[-1]}",
			data: map[string]interface{}{
				"path": "/home/user/projects/test",
			},
			expected: "Dir: test",
		},
		{
			name:     "missing field",
			template: "Value: {.missing}",
			data:     map[string]interface{}{},
			expected: "Value: ",
		},
		{
			name:     "multiple templates",
			template: "{.name} - {.version}",
			data: map[string]interface{}{
				"name":    "ccstatusline",
				"version": "1.0.0",
			},
			expected: "ccstatusline - 1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processTemplate(tt.template, tt.data)
			if result != tt.expected {
				t.Errorf("processTemplate() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestJQValueToString(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "string value",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "nil value",
			value:    nil,
			expected: "",
		},
		{
			name:     "bool true",
			value:    true,
			expected: "true",
		},
		{
			name:     "bool false",
			value:    false,
			expected: "false",
		},
		{
			name:     "number",
			value:    42,
			expected: "42",
		},
		{
			name:     "float",
			value:    3.14,
			expected: "3.14",
		},
		{
			name:     "array",
			value:    []interface{}{"a", "b", "c"},
			expected: `["a","b","c"]`,
		},
		{
			name: "object",
			value: map[string]interface{}{
				"key": "value",
			},
			expected: `{"key":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := jqValueToString(tt.value)
			if result != tt.expected {
				t.Errorf("jqValueToString() = %q, want %q", result, tt.expected)
			}
		})
	}
}
