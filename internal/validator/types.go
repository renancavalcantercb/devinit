package validator

import "github.com/renan-dev/devinit/internal/template"

// ValidationLevel defines how strict validation should be
type ValidationLevel int

const (
	// ValidationNone skips all validation (--no-validate flag)
	ValidationNone ValidationLevel = iota
	// ValidationBasic performs basic checks (default)
	ValidationBasic
	// ValidationStrict performs strict version checking (--strict flag)
	ValidationStrict
)

// String returns the string representation of ValidationLevel
func (v ValidationLevel) String() string {
	switch v {
	case ValidationNone:
		return "none"
	case ValidationBasic:
		return "basic"
	case ValidationStrict:
		return "strict"
	default:
		return "unknown"
	}
}

// ValidationResult contains the results of validation
type ValidationResult struct {
	Errors   []ValidationError
	Warnings []ValidationError
}

// HasErrors returns true if there are any errors
func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// HasWarnings returns true if there are any warnings
func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// ValidationError represents a validation error or warning
type ValidationError struct {
	Command     string
	Message     string
	InstallHint string
}

// Error implements the error interface
func (e ValidationError) Error() string {
	return e.Message
}

// Requirement represents a system requirement
type Requirement struct {
	Command     string
	Version     string
	Required    bool
	When        string
	InstallHint string
}

// FromTemplateRequirement converts a template.SystemRequirement to a Requirement
func FromTemplateRequirement(tr template.SystemRequirement) Requirement {
	return Requirement{
		Command:     tr.Command,
		Version:     tr.Version,
		Required:    tr.Required,
		When:        tr.When,
		InstallHint: tr.InstallHint,
	}
}
