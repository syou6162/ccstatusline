package main

import (
	"fmt"
	"strings"
)

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

var bgColorMap = map[string]string{
	// Normal background colors
	"bg_black":   "\033[40m",
	"bg_red":     "\033[41m",
	"bg_green":   "\033[42m",
	"bg_yellow":  "\033[43m",
	"bg_blue":    "\033[44m",
	"bg_magenta": "\033[45m",
	"bg_cyan":    "\033[46m",
	"bg_white":   "\033[47m",

	// Bright background colors
	"bg_gray":           "\033[100m", // bright black
	"bg_bright_red":     "\033[101m",
	"bg_bright_green":   "\033[102m",
	"bg_bright_yellow":  "\033[103m",
	"bg_bright_blue":    "\033[104m",
	"bg_bright_magenta": "\033[105m",
	"bg_bright_cyan":    "\033[106m",
	"bg_bright_white":   "\033[107m",
}

// applyColor applies ANSI color codes to text
func applyColor(text, color string) string {
	if color == "" {
		return text
	}

	var colorCode string

	// Check if it's a background color (starts with bg_)
	if strings.HasPrefix(color, "bg_") {
		colorCode = bgColorMap[color]
	} else {
		// It's a foreground color
		colorCode = colorMap[color]
	}

	if colorCode == "" {
		// Unknown color, return text as-is
		return text
	}

	return fmt.Sprintf("%s%s%s", colorCode, text, resetCode)
}
