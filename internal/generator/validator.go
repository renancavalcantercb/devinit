package generator

import (
	"fmt"
	"os"
	"regexp"
)

var projectNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// ValidateProjectName validates a project name for security and correctness
//
// Security checks:
// - Prevents path traversal attacks (../, absolute paths)
// - Ensures safe filesystem operations
//
// Format requirements:
// - Must start with lowercase letter
// - Only lowercase letters, numbers, and hyphens allowed
// - This ensures compatibility across filesystems and platforms
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check for path traversal attempts
	if name == "." || name == ".." {
		return fmt.Errorf("invalid project name: '.' and '..' are not allowed")
	}

	// Check for path separators (security: prevent path traversal)
	for _, char := range name {
		if char == '/' || char == '\\' {
			return fmt.Errorf("invalid project name: path separators are not allowed")
		}
	}

	// Validate against pattern
	if !projectNamePattern.MatchString(name) {
		return fmt.Errorf("invalid project name: must start with lowercase letter and contain only lowercase letters, numbers, and hyphens")
	}

	// Check if directory already exists
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("directory '%s' already exists", name)
	}

	return nil
}
