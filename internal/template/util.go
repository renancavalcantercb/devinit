package template

import (
	"regexp"
	"strings"
	"unicode"
)

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	// Replace hyphens with underscores
	s = strings.ReplaceAll(s, "-", "_")

	// Insert underscore before uppercase letters
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}

	// Clean up multiple underscores
	re := regexp.MustCompile(`_+`)
	return re.ReplaceAllString(result.String(), "_")
}

// toCamelCase converts a string to camelCase
func toCamelCase(s string) string {
	// Split on hyphens, underscores, and spaces
	parts := regexp.MustCompile(`[-_\s]+`).Split(s, -1)

	if len(parts) == 0 {
		return ""
	}

	var result strings.Builder
	for i, part := range parts {
		if part == "" {
			continue
		}

		if i == 0 {
			// First part is lowercase
			result.WriteString(strings.ToLower(part))
		} else {
			// Capitalize first letter of subsequent parts
			result.WriteString(capitalize(part))
		}
	}

	return result.String()
}

// toPascalCase converts a string to PascalCase
func toPascalCase(s string) string {
	// Split on hyphens, underscores, and spaces
	parts := regexp.MustCompile(`[-_\s]+`).Split(s, -1)

	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		result.WriteString(capitalize(part))
	}

	return result.String()
}

// toKebabCase converts a string to kebab-case
func toKebabCase(s string) string {
	// Replace underscores with hyphens
	s = strings.ReplaceAll(s, "_", "-")

	// Insert hyphen before uppercase letters
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(r))
	}

	// Clean up multiple hyphens
	re := regexp.MustCompile(`-+`)
	return re.ReplaceAllString(result.String(), "-")
}

// capitalize capitalizes the first letter of a string
func capitalize(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
