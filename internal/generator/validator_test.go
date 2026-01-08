package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		errorMsg  string
	}{
		// Valid names
		{
			name:      "valid simple name",
			input:     "myproject",
			wantError: false,
		},
		{
			name:      "valid name with hyphens",
			input:     "my-api-service",
			wantError: false,
		},
		{
			name:      "valid name with numbers",
			input:     "api2024",
			wantError: false,
		},
		{
			name:      "valid name with hyphens and numbers",
			input:     "my-api-v2",
			wantError: false,
		},

		// Invalid: empty name
		{
			name:      "empty name",
			input:     "",
			wantError: true,
			errorMsg:  "project name cannot be empty",
		},

		// Invalid: uppercase letters
		{
			name:      "uppercase letters",
			input:     "MyProject",
			wantError: true,
			errorMsg:  "must start with lowercase letter",
		},
		{
			name:      "all uppercase",
			input:     "PROJECT",
			wantError: true,
			errorMsg:  "must start with lowercase letter",
		},

		// Invalid: starts with number
		{
			name:      "starts with number",
			input:     "2project",
			wantError: true,
			errorMsg:  "must start with lowercase letter",
		},

		// Invalid: starts with hyphen
		{
			name:      "starts with hyphen",
			input:     "-project",
			wantError: true,
			errorMsg:  "must start with lowercase letter",
		},

		// Security: path traversal attempts
		{
			name:      "path traversal with ../",
			input:     "../project",
			wantError: true,
			errorMsg:  "path separators are not allowed",
		},
		{
			name:      "path traversal parent dir",
			input:     "..",
			wantError: true,
			errorMsg:  "'.' and '..' are not allowed",
		},
		{
			name:      "current directory",
			input:     ".",
			wantError: true,
			errorMsg:  "'.' and '..' are not allowed",
		},
		{
			name:      "absolute path",
			input:     "/tmp/project",
			wantError: true,
			errorMsg:  "path separators are not allowed",
		},
		{
			name:      "windows path",
			input:     "C:\\project",
			wantError: true,
			errorMsg:  "path separators are not allowed",
		},
		{
			name:      "relative path with forward slash",
			input:     "dir/project",
			wantError: true,
			errorMsg:  "path separators are not allowed",
		},
		{
			name:      "relative path with backslash",
			input:     "dir\\project",
			wantError: true,
			errorMsg:  "path separators are not allowed",
		},

		// Invalid: special characters
		{
			name:      "underscore",
			input:     "my_project",
			wantError: true,
			errorMsg:  "must start with lowercase letter",
		},
		{
			name:      "space",
			input:     "my project",
			wantError: true,
			errorMsg:  "must start with lowercase letter",
		},
		{
			name:      "special characters",
			input:     "my@project",
			wantError: true,
			errorMsg:  "must start with lowercase letter",
		},
		{
			name:      "dots",
			input:     "my.project",
			wantError: true,
			errorMsg:  "must start with lowercase letter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateProjectName(%q) expected error, got nil", tt.input)
					return
				}
				if tt.errorMsg != "" && !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("ValidateProjectName(%q) error = %v, want error containing %q", tt.input, err, tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateProjectName(%q) unexpected error: %v", tt.input, err)
				}
			}
		})
	}
}

func TestValidateProjectName_ExistingDirectory(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	existingDir := filepath.Join(tmpDir, "existing-project")

	if err := os.Mkdir(existingDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Change to temp directory for testing
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test validation against existing directory
	err = ValidateProjectName("existing-project")
	if err == nil {
		t.Error("ValidateProjectName should fail for existing directory")
	}
	if !containsString(err.Error(), "already exists") {
		t.Errorf("Error message should mention 'already exists', got: %v", err)
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
