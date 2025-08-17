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
	name = normalizeColorName(name)

	// Check if it's a background color
	if strings.HasPrefix(name, "bg_") {
		bgColor := strings.TrimPrefix(name, "bg_")

		// Handle bright background colors
		if strings.HasPrefix(bgColor, "bright_") || strings.HasPrefix(bgColor, "light") {
			bgColor = normalizeBrightName(bgColor)
			if c, ok := color.ExBgColors[bgColor]; ok {
				return color.New(c)
			}
		}

		// Normal background colors
		if c, ok := color.BgColors[bgColor]; ok {
			return color.New(c)
		}

		// Special case for gray/grey (treated as darkGray in background)
		if bgColor == "gray" || bgColor == "grey" {
			if c, ok := color.ExBgColors["darkGray"]; ok {
				return color.New(c)
			}
		}
	}

	// Handle bright/light foreground colors
	if strings.HasPrefix(name, "bright_") || strings.HasPrefix(name, "light") {
		fgColor := normalizeBrightName(name)
		if c, ok := color.ExFgColors[fgColor]; ok {
			return color.New(c)
		}
	}

	// Check if it's a style option (bold, italic, etc.)
	if c, ok := color.AllOptions[name]; ok {
		return color.New(c)
	}

	// Normal foreground colors
	if c, ok := color.FgColors[name]; ok {
		return color.New(c)
	}

	// Special case for gray/grey (treated as darkGray)
	if name == "gray" || name == "grey" {
		if c, ok := color.ExFgColors["darkGray"]; ok {
			return color.New(c)
		}
	}

	return nil
}

// normalizeColorName handles common aliases and variations
func normalizeColorName(name string) string {
	// Handle grey -> gray for consistency
	name = strings.ReplaceAll(name, "grey", "gray")
	// Handle dim -> fuzzy (gookit/color uses "fuzzy" for dim)
	if name == "dim" {
		return "fuzzy"
	}
	// Handle underline -> underscore (gookit/color uses "underscore")
	if name == "underline" {
		return "underscore"
	}
	// Handle strike -> strikethrough (not in AllOptions, needs special handling)
	if name == "strike" || name == "strikethrough" {
		// Note: gookit/color doesn't have strikethrough in AllOptions
		// We'll handle this case in parseColorName
		return name
	}
	return name
}

// normalizeBrightName converts our naming convention to gookit/color's convention
func normalizeBrightName(name string) string {
	// Convert bright_* to light* (gookit/color uses "light" prefix)
	name = strings.ReplaceAll(name, "bright_", "light")
	name = strings.ReplaceAll(name, "_", "")

	// Capitalize the color name after "light"
	if strings.HasPrefix(name, "light") && len(name) > 5 {
		colorPart := name[5:]
		name = "light" + strings.Title(colorPart)
	}

	return name
}
