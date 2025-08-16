package main

import "fmt"

// ANSI color codes
const (
	resetCode = "\033[0m"
)

var colorMap = map[string]string{
	// Normal colors
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",

	// Bright colors
	"gray":           "\033[90m", // bright black
	"bright_red":     "\033[91m",
	"bright_green":   "\033[92m",
	"bright_yellow":  "\033[93m",
	"bright_blue":    "\033[94m",
	"bright_magenta": "\033[95m",
	"bright_cyan":    "\033[96m",
	"bright_white":   "\033[97m",
}

// applyColor applies ANSI color codes to text
func applyColor(text, color string) string {
	if color == "" {
		return text
	}

	colorCode, ok := colorMap[color]
	if !ok {
		// Unknown color, return text as-is
		return text
	}

	return fmt.Sprintf("%s%s%s", colorCode, text, resetCode)
}
