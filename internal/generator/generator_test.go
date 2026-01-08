package generator

import (
	"testing"

	"github.com/renan-dev/devinit/internal/template"
)

func TestEvaluateCondition(t *testing.T) {
	gen := &Generator{}

	// Create test context
	variables := map[string]interface{}{
		"IncludeDocker": true,
		"IncludeTests":  false,
		"CustomFlag":    true,
	}
	ctx := template.NewContext("test-project", "/tmp/test", variables, &template.Template{})

	tests := []struct {
		name      string
		condition string
		want      bool
	}{
		// Template syntax with braces
		{
			name:      "template syntax with braces - true",
			condition: "{{ .IncludeDocker }}",
			want:      true,
		},
		{
			name:      "template syntax with braces - false",
			condition: "{{ .IncludeTests }}",
			want:      false,
		},

		// Without braces, with dot
		{
			name:      "with dot prefix - true",
			condition: ".IncludeDocker",
			want:      true,
		},
		{
			name:      "with dot prefix - false",
			condition: ".IncludeTests",
			want:      false,
		},

		// Plain variable name
		{
			name:      "plain variable name - true",
			condition: "IncludeDocker",
			want:      true,
		},
		{
			name:      "plain variable name - false",
			condition: "IncludeTests",
			want:      false,
		},

		// With whitespace
		{
			name:      "with whitespace",
			condition: "  {{ .IncludeDocker }}  ",
			want:      true,
		},
		{
			name:      "with inner whitespace",
			condition: "{{  .IncludeDocker  }}",
			want:      true,
		},

		// Custom variables via map
		{
			name:      "custom variable - true",
			condition: "CustomFlag",
			want:      true,
		},

		// Non-existent variable (should be false)
		{
			name:      "non-existent variable",
			condition: "NonExistent",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.evaluateCondition(tt.condition, ctx)
			if got != tt.want {
				t.Errorf("evaluateCondition(%q) = %v, want %v", tt.condition, got, tt.want)
			}
		})
	}
}

func TestShouldGenerateFile(t *testing.T) {
	gen := &Generator{}

	// Create test context
	variables := map[string]interface{}{
		"IncludeDocker": true,
		"IncludeTests":  false,
		"Database":      "postgres",
	}
	ctx := template.NewContext("test-project", "/tmp/test", variables, &template.Template{})

	tests := []struct {
		name     string
		fileSpec template.FileSpec
		want     bool
	}{
		{
			name: "no conditions - always generate",
			fileSpec: template.FileSpec{
				Source:      "main.py.tmpl",
				Destination: "src/main.py",
				Conditions:  []string{},
			},
			want: true,
		},
		{
			name: "single condition - true",
			fileSpec: template.FileSpec{
				Source:      "Dockerfile",
				Destination: "Dockerfile",
				Conditions:  []string{"{{ .IncludeDocker }}"},
			},
			want: true,
		},
		{
			name: "single condition - false",
			fileSpec: template.FileSpec{
				Source:      "test_main.py.tmpl",
				Destination: "tests/test_main.py",
				Conditions:  []string{"{{ .IncludeTests }}"},
			},
			want: false,
		},
		{
			name: "multiple conditions - all true",
			fileSpec: template.FileSpec{
				Source:      "docker-compose.yml.tmpl",
				Destination: "docker-compose.yml",
				Conditions:  []string{"{{ .IncludeDocker }}", "IncludeDocker"},
			},
			want: true,
		},
		{
			name: "multiple conditions - one false",
			fileSpec: template.FileSpec{
				Source:      "special-file.txt",
				Destination: "special-file.txt",
				Conditions:  []string{"{{ .IncludeDocker }}", "{{ .IncludeTests }}"},
			},
			want: false,
		},
		{
			name: "multiple conditions - all false",
			fileSpec: template.FileSpec{
				Source:      "never-generated.txt",
				Destination: "never-generated.txt",
				Conditions:  []string{"{{ .IncludeTests }}", "NonExistent"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.shouldGenerateFile(tt.fileSpec, ctx)
			if got != tt.want {
				t.Errorf("shouldGenerateFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
