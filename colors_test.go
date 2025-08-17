package main

import (
	"testing"
)

func TestApplyColor(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		color    string
		expected string
	}{
		{
			name:     "cyan color",
			text:     "Hello",
			color:    "cyan",
			expected: "\033[36mHello\033[0m",
		},
		{
			name:     "red color",
			text:     "Error",
			color:    "red",
			expected: "\033[31mError\033[0m",
		},
		{
			name:     "green color",
			text:     "Success",
			color:    "green",
			expected: "\033[32mSuccess\033[0m",
		},
		{
			name:     "gray color",
			text:     "Debug",
			color:    "gray",
			expected: "\033[90mDebug\033[0m",
		},
		{
			name:     "bright_blue color",
			text:     "Info",
			color:    "bright_blue",
			expected: "\033[94mInfo\033[0m",
		},
		{
			name:     "empty color",
			text:     "No color",
			color:    "",
			expected: "No color",
		},
		{
			name:     "unknown color",
			text:     "Unknown",
			color:    "invalid",
			expected: "Unknown",
		},
		{
			name:     "bg_magenta background",
			text:     "Background",
			color:    "bg_magenta",
			expected: "\033[45mBackground\033[0m",
		},
		{
			name:     "bg_yellow background",
			text:     "Warning",
			color:    "bg_yellow",
			expected: "\033[43mWarning\033[0m",
		},
		{
			name:     "bg_bright_cyan background",
			text:     "Info",
			color:    "bg_bright_cyan",
			expected: "\033[106mInfo\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyColor(tt.text, tt.color)
			if result != tt.expected {
				t.Errorf("applyColor(%q, %q) = %q, want %q", tt.text, tt.color, result, tt.expected)
			}
		})
	}
}
