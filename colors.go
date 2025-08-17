package main

import (
	"github.com/gookit/color"
	"strings"
)

// applyColor applies color to text using gookit/color library
func applyColor(text, colorName string) string {
	if text == "" || colorName == "" {
		return text
	}

	// Parse color name and create a Style using the library
	style := parseColorName(colorName)
	if style == nil {
		return text // Unknown color
	}

	// Use the Style's Sprint method to render the text
	return style.Sprint(text)
}

// parseColorName converts our color naming convention to gookit/color Style
func parseColorName(name string) color.Style {
	// Handle aliases
	name = strings.ReplaceAll(name, "grey", "gray")

	// Check if it's a background color
	if strings.HasPrefix(name, "bg_") {
		bgColor := strings.TrimPrefix(name, "bg_")

		// Handle bright background colors
		if strings.HasPrefix(bgColor, "bright_") {
			bgColor = strings.TrimPrefix(bgColor, "bright_")
			switch bgColor {
			case "black", "gray":
				return color.New(color.BgGray)
			case "red":
				return color.New(color.BgLightRed)
			case "green":
				return color.New(color.BgLightGreen)
			case "yellow":
				return color.New(color.BgLightYellow)
			case "blue":
				return color.New(color.BgLightBlue)
			case "magenta":
				return color.New(color.BgLightMagenta)
			case "cyan":
				return color.New(color.BgLightCyan)
			case "white":
				return color.New(color.BgLightWhite)
			}
		}

		// Normal background colors
		switch bgColor {
		case "black":
			return color.New(color.BgBlack)
		case "red":
			return color.New(color.BgRed)
		case "green":
			return color.New(color.BgGreen)
		case "yellow":
			return color.New(color.BgYellow)
		case "blue":
			return color.New(color.BgBlue)
		case "magenta":
			return color.New(color.BgMagenta)
		case "cyan":
			return color.New(color.BgCyan)
		case "white":
			return color.New(color.BgWhite)
		case "gray":
			return color.New(color.BgGray)
		}
	}

	// Handle bright foreground colors
	if strings.HasPrefix(name, "bright_") {
		fgColor := strings.TrimPrefix(name, "bright_")
		switch fgColor {
		case "black", "gray":
			return color.New(color.FgGray)
		case "red":
			return color.New(color.FgLightRed)
		case "green":
			return color.New(color.FgLightGreen)
		case "yellow":
			return color.New(color.FgLightYellow)
		case "blue":
			return color.New(color.FgLightBlue)
		case "magenta":
			return color.New(color.FgLightMagenta)
		case "cyan":
			return color.New(color.FgLightCyan)
		case "white":
			return color.New(color.FgLightWhite)
		}
	}

	// Normal foreground colors and styles
	switch name {
	case "black":
		return color.New(color.FgBlack)
	case "red":
		return color.New(color.FgRed)
	case "green":
		return color.New(color.FgGreen)
	case "yellow":
		return color.New(color.FgYellow)
	case "blue":
		return color.New(color.FgBlue)
	case "magenta":
		return color.New(color.FgMagenta)
	case "cyan":
		return color.New(color.FgCyan)
	case "white":
		return color.New(color.FgWhite)
	case "gray":
		return color.New(color.FgGray)
	// Styles
	case "bold":
		return color.New(color.OpBold)
	case "dim":
		return color.New(color.OpFuzzy)
	case "italic":
		return color.New(color.OpItalic)
	case "underline":
		return color.New(color.OpUnderscore)
	case "blink":
		return color.New(color.OpBlink)
	case "reverse":
		return color.New(color.OpReverse)
	case "hidden":
		return color.New(color.OpConcealed)
	case "strike":
		return color.New(color.OpStrikethrough)
	}

	return nil
}
