package validator

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// SystemValidator validates system requirements
type SystemValidator struct {
	Level ValidationLevel
}

// NewSystemValidator creates a new system validator
func NewSystemValidator(level ValidationLevel) *SystemValidator {
	return &SystemValidator{
		Level: level,
	}
}

// Validate checks if all requirements are met
func (v *SystemValidator) Validate(reqs []Requirement) (*ValidationResult, error) {
	result := &ValidationResult{
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}

	for _, req := range reqs {
		// TODO: Evaluate When condition when template context is available
		// For now, we check all requirements

		exists, version, err := v.CheckCommand(req.Command)

		if err != nil {
			// Error checking command
			valErr := ValidationError{
				Command:     req.Command,
				Message:     fmt.Sprintf("error checking %s: %v", req.Command, err),
				InstallHint: req.InstallHint,
			}

			if req.Required {
				result.Errors = append(result.Errors, valErr)
			} else {
				result.Warnings = append(result.Warnings, valErr)
			}
			continue
		}

		if !exists {
			// Command not found
			valErr := ValidationError{
				Command:     req.Command,
				Message:     fmt.Sprintf("%s not found", req.Command),
				InstallHint: req.InstallHint,
			}

			if req.Required {
				result.Errors = append(result.Errors, valErr)
			} else {
				result.Warnings = append(result.Warnings, valErr)
			}
			continue
		}

		// Check version if specified
		if req.Version != "" && version != "" {
			matches, err := v.CompareVersion(version, req.Version)
			if err != nil {
				valErr := ValidationError{
					Command:     req.Command,
					Message:     fmt.Sprintf("error comparing %s version: %v", req.Command, err),
					InstallHint: req.InstallHint,
				}

				if v.Level == ValidationStrict {
					result.Errors = append(result.Errors, valErr)
				} else {
					result.Warnings = append(result.Warnings, valErr)
				}
				continue
			}

			if !matches {
				valErr := ValidationError{
					Command: req.Command,
					Message: fmt.Sprintf("%s version %s does not match requirement %s",
						req.Command, version, req.Version),
					InstallHint: req.InstallHint,
				}

				if v.Level == ValidationStrict {
					result.Errors = append(result.Errors, valErr)
				} else {
					result.Warnings = append(result.Warnings, valErr)
				}
			}
		}
	}

	return result, nil
}

// CheckCommand checks if a command exists and returns its version
func (v *SystemValidator) CheckCommand(cmd string) (exists bool, version string, err error) {
	// Check if command exists using 'which' on Unix or 'where' on Windows
	_, err = exec.LookPath(cmd)
	if err != nil {
		return false, "", nil
	}

	// Try to get version
	version, _ = v.getCommandVersion(cmd)

	return true, version, nil
}

// getCommandVersion attempts to get the version of a command
func (v *SystemValidator) getCommandVersion(cmd string) (string, error) {
	// Common version flags
	versionFlags := []string{"--version", "-version", "-v", "version"}

	for _, flag := range versionFlags {
		output, err := exec.Command(cmd, flag).CombinedOutput()
		if err != nil {
			continue
		}

		// Try to extract version number
		version := extractVersion(string(output))
		if version != "" {
			return version, nil
		}
	}

	return "", fmt.Errorf("unable to determine version")
}

// extractVersion extracts a semantic version from command output
func extractVersion(output string) string {
	// Pattern for semantic versioning (e.g., 3.11.5, v1.2.3, 20.0.1)
	patterns := []string{
		`v?(\d+\.\d+\.\d+)`,           // Standard semver
		`v?(\d+\.\d+)`,                 // Major.minor
		`version\s+v?(\d+\.\d+\.\d+)`, // With "version" prefix
		`(\d+\.\d+\.\d+)`,              // Just numbers
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(output)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// CompareVersion compares a version string against a requirement
// Supports: >=, >, <=, <, =, ^, ~
func (v *SystemValidator) CompareVersion(current, requirement string) (bool, error) {
	requirement = strings.TrimSpace(requirement)

	// Parse operator and version
	operator := ""
	requiredVersion := requirement

	if strings.HasPrefix(requirement, ">=") {
		operator = ">="
		requiredVersion = strings.TrimSpace(requirement[2:])
	} else if strings.HasPrefix(requirement, "<=") {
		operator = "<="
		requiredVersion = strings.TrimSpace(requirement[2:])
	} else if strings.HasPrefix(requirement, ">") {
		operator = ">"
		requiredVersion = strings.TrimSpace(requirement[1:])
	} else if strings.HasPrefix(requirement, "<") {
		operator = "<"
		requiredVersion = strings.TrimSpace(requirement[1:])
	} else if strings.HasPrefix(requirement, "=") {
		operator = "="
		requiredVersion = strings.TrimSpace(requirement[1:])
	} else if strings.HasPrefix(requirement, "^") {
		// Caret: ^1.2.3 allows >=1.2.3 but <2.0.0
		operator = "^"
		requiredVersion = strings.TrimSpace(requirement[1:])
	} else if strings.HasPrefix(requirement, "~") {
		// Tilde: ~1.2.3 allows >=1.2.3 but <1.3.0
		operator = "~"
		requiredVersion = strings.TrimSpace(requirement[1:])
	} else {
		// No operator means exact match
		operator = "="
	}

	// Parse versions
	currentParts, err := parseVersion(current)
	if err != nil {
		return false, fmt.Errorf("invalid current version %s: %w", current, err)
	}

	requiredParts, err := parseVersion(requiredVersion)
	if err != nil {
		return false, fmt.Errorf("invalid required version %s: %w", requiredVersion, err)
	}

	// Compare based on operator
	comparison := compareVersionParts(currentParts, requiredParts)

	switch operator {
	case ">=":
		return comparison >= 0, nil
	case ">":
		return comparison > 0, nil
	case "<=":
		return comparison <= 0, nil
	case "<":
		return comparison < 0, nil
	case "=":
		return comparison == 0, nil
	case "^":
		// ^1.2.3 allows >=1.2.3 but <2.0.0
		if comparison < 0 {
			return false, nil
		}
		// Check if major version is the same
		return currentParts[0] == requiredParts[0], nil
	case "~":
		// ~1.2.3 allows >=1.2.3 but <1.3.0
		if comparison < 0 {
			return false, nil
		}
		// Check if major and minor versions are the same
		return currentParts[0] == requiredParts[0] && currentParts[1] == requiredParts[1], nil
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}

// parseVersion parses a version string into [major, minor, patch]
func parseVersion(version string) ([3]int, error) {
	version = strings.TrimPrefix(version, "v")
	parts := strings.Split(version, ".")

	var result [3]int
	for i := 0; i < 3; i++ {
		if i < len(parts) {
			num, err := strconv.Atoi(parts[i])
			if err != nil {
				return result, fmt.Errorf("invalid version component %s: %w", parts[i], err)
			}
			result[i] = num
		} else {
			result[i] = 0
		}
	}

	return result, nil
}

// compareVersionParts compares two version arrays
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersionParts(v1, v2 [3]int) int {
	for i := 0; i < 3; i++ {
		if v1[i] < v2[i] {
			return -1
		}
		if v1[i] > v2[i] {
			return 1
		}
	}
	return 0
}
