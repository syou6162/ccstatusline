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
		// Foreground colors
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
		// Background colors
		{
			name:     "bg_blue color",
			text:     "Background",
			color:    "bg_blue",
			expected: "\033[44mBackground\033[0m",
		},
		{
			name:     "bg_bright_yellow color",
			text:     "Warning",
			color:    "bg_bright_yellow",
			expected: "\033[103mWarning\033[0m",
		},
		// Style modifiers
		{
			name:     "bold style",
			text:     "Bold",
			color:    "bold",
			expected: "\033[1mBold\033[0m",
		},
		{
			name:     "underline style",
			text:     "Underlined",
			color:    "underline",
			expected: "\033[4mUnderlined\033[0m",
		},
		// Edge cases
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
		// Aliases
		{
			name:     "grey alias",
			text:     "Grey",
			color:    "grey",
			expected: "\033[90mGrey\033[0m",
		},
		{
			name:     "bg_grey alias",
			text:     "BG Grey",
			color:    "bg_grey",
			expected: "\033[100mBG Grey\033[0m",
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
