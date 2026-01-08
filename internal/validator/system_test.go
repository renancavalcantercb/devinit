package validator

import (
	"testing"
)

func TestCompareVersion(t *testing.T) {
	validator := NewSystemValidator(ValidationBasic)

	tests := []struct {
		name        string
		current     string
		requirement string
		want        bool
		wantErr     bool
	}{
		// Exact match
		{
			name:        "exact match",
			current:     "3.11.5",
			requirement: "=3.11.5",
			want:        true,
		},
		{
			name:        "exact match without operator",
			current:     "3.11.5",
			requirement: "3.11.5",
			want:        true,
		},
		{
			name:        "exact match fail",
			current:     "3.11.5",
			requirement: "=3.11.6",
			want:        false,
		},

		// Greater than or equal
		{
			name:        "greater than or equal - equal",
			current:     "3.11.5",
			requirement: ">=3.11.5",
			want:        true,
		},
		{
			name:        "greater than or equal - greater",
			current:     "3.12.0",
			requirement: ">=3.11.5",
			want:        true,
		},
		{
			name:        "greater than or equal - less",
			current:     "3.10.0",
			requirement: ">=3.11.5",
			want:        false,
		},

		// Greater than
		{
			name:        "greater than - true",
			current:     "3.12.0",
			requirement: ">3.11.5",
			want:        true,
		},
		{
			name:        "greater than - equal",
			current:     "3.11.5",
			requirement: ">3.11.5",
			want:        false,
		},
		{
			name:        "greater than - less",
			current:     "3.10.0",
			requirement: ">3.11.5",
			want:        false,
		},

		// Less than or equal
		{
			name:        "less than or equal - equal",
			current:     "3.11.5",
			requirement: "<=3.11.5",
			want:        true,
		},
		{
			name:        "less than or equal - less",
			current:     "3.10.0",
			requirement: "<=3.11.5",
			want:        true,
		},
		{
			name:        "less than or equal - greater",
			current:     "3.12.0",
			requirement: "<=3.11.5",
			want:        false,
		},

		// Less than
		{
			name:        "less than - true",
			current:     "3.10.0",
			requirement: "<3.11.5",
			want:        true,
		},
		{
			name:        "less than - equal",
			current:     "3.11.5",
			requirement: "<3.11.5",
			want:        false,
		},
		{
			name:        "less than - greater",
			current:     "3.12.0",
			requirement: "<3.11.5",
			want:        false,
		},

		// Caret (^) - allows patch and minor updates
		{
			name:        "caret - same version",
			current:     "1.2.3",
			requirement: "^1.2.3",
			want:        true,
		},
		{
			name:        "caret - patch update",
			current:     "1.2.4",
			requirement: "^1.2.3",
			want:        true,
		},
		{
			name:        "caret - minor update",
			current:     "1.3.0",
			requirement: "^1.2.3",
			want:        true,
		},
		{
			name:        "caret - major update (should fail)",
			current:     "2.0.0",
			requirement: "^1.2.3",
			want:        false,
		},
		{
			name:        "caret - older version (should fail)",
			current:     "1.2.2",
			requirement: "^1.2.3",
			want:        false,
		},

		// Tilde (~) - allows only patch updates
		{
			name:        "tilde - same version",
			current:     "1.2.3",
			requirement: "~1.2.3",
			want:        true,
		},
		{
			name:        "tilde - patch update",
			current:     "1.2.4",
			requirement: "~1.2.3",
			want:        true,
		},
		{
			name:        "tilde - minor update (should fail)",
			current:     "1.3.0",
			requirement: "~1.2.3",
			want:        false,
		},
		{
			name:        "tilde - major update (should fail)",
			current:     "2.0.0",
			requirement: "~1.2.3",
			want:        false,
		},

		// Version with 'v' prefix
		{
			name:        "version with v prefix",
			current:     "v3.11.5",
			requirement: ">=3.11.0",
			want:        true,
		},
		{
			name:        "requirement with v prefix",
			current:     "3.11.5",
			requirement: ">=v3.11.0",
			want:        true,
		},

		// Major.minor versions (no patch)
		{
			name:        "major.minor version",
			current:     "3.11",
			requirement: ">=3.10",
			want:        true,
		},
		{
			name:        "major.minor version exact",
			current:     "20.0",
			requirement: "=20.0",
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validator.CompareVersion(tt.current, tt.requirement)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareVersion(%s, %s) = %v, want %v", tt.current, tt.requirement, got, tt.want)
			}
		})
	}
}

func TestCheckCommand(t *testing.T) {
	validator := NewSystemValidator(ValidationBasic)

	tests := []struct {
		name          string
		command       string
		wantExists    bool
		wantVersion   bool // true if we expect to get some version
	}{
		{
			name:        "existing command - go",
			command:     "go",
			wantExists:  true,
			wantVersion: true,
		},
		{
			name:        "existing command - git",
			command:     "git",
			wantExists:  true,
			wantVersion: true,
		},
		{
			name:       "non-existing command",
			command:    "this-command-definitely-does-not-exist-12345",
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, version, err := validator.CheckCommand(tt.command)

			if err != nil {
				t.Errorf("CheckCommand() unexpected error: %v", err)
			}

			if exists != tt.wantExists {
				t.Errorf("CheckCommand(%s) exists = %v, want %v", tt.command, exists, tt.wantExists)
			}

			if tt.wantVersion && version == "" && exists {
				t.Logf("Warning: Could not determine version for %s (this may be okay)", tt.command)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name         string
		level        ValidationLevel
		requirements []Requirement
		wantErrors   int
		wantWarnings int
	}{
		{
			name:  "all requirements met",
			level: ValidationBasic,
			requirements: []Requirement{
				{
					Command:  "go",
					Required: true,
				},
			},
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name:  "required command missing",
			level: ValidationBasic,
			requirements: []Requirement{
				{
					Command:  "this-does-not-exist",
					Required: true,
				},
			},
			wantErrors:   1,
			wantWarnings: 0,
		},
		{
			name:  "optional command missing",
			level: ValidationBasic,
			requirements: []Requirement{
				{
					Command:  "this-does-not-exist",
					Required: false,
				},
			},
			wantErrors:   0,
			wantWarnings: 1,
		},
		{
			name:  "mixed requirements",
			level: ValidationBasic,
			requirements: []Requirement{
				{
					Command:  "go",
					Required: true,
				},
				{
					Command:  "this-does-not-exist",
					Required: false,
				},
			},
			wantErrors:   0,
			wantWarnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewSystemValidator(tt.level)
			result, err := validator.Validate(tt.requirements)

			if err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
				return
			}

			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Validate() errors = %d, want %d", len(result.Errors), tt.wantErrors)
				for _, e := range result.Errors {
					t.Logf("  Error: %v", e)
				}
			}

			if len(result.Warnings) != tt.wantWarnings {
				t.Errorf("Validate() warnings = %d, want %d", len(result.Warnings), tt.wantWarnings)
				for _, w := range result.Warnings {
					t.Logf("  Warning: %v", w)
				}
			}
		})
	}
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{
			name:   "python version",
			output: "Python 3.11.5",
			want:   "3.11.5",
		},
		{
			name:   "go version",
			output: "go version go1.21.0 darwin/arm64",
			want:   "1.21.0",
		},
		{
			name:   "node version",
			output: "v20.9.0",
			want:   "20.9.0",
		},
		{
			name:   "docker version",
			output: "Docker version 24.0.6, build ed223bc",
			want:   "24.0.6",
		},
		{
			name:   "version with v prefix",
			output: "version v1.2.3",
			want:   "1.2.3",
		},
		{
			name:   "major.minor version",
			output: "Version 20.0",
			want:   "20.0",
		},
		{
			name:   "no version found",
			output: "some output without version",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractVersion(tt.output)
			if got != tt.want {
				t.Errorf("extractVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    [3]int
		wantErr bool
	}{
		{
			name:    "full semver",
			version: "3.11.5",
			want:    [3]int{3, 11, 5},
		},
		{
			name:    "with v prefix",
			version: "v1.2.3",
			want:    [3]int{1, 2, 3},
		},
		{
			name:    "major.minor only",
			version: "20.0",
			want:    [3]int{20, 0, 0},
		},
		{
			name:    "major only",
			version: "3",
			want:    [3]int{3, 0, 0},
		},
		{
			name:    "invalid version",
			version: "not.a.version",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseVersion(%s) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
