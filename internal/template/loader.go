package template

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader loads templates from the filesystem
type Loader struct {
	templatesDir string
}

// NewLoader creates a new template loader
func NewLoader(templatesDir string) *Loader {
	return &Loader{
		templatesDir: templatesDir,
	}
}

// Load loads a template by name (e.g., "python/fastapi")
func (l *Loader) Load(name string) (*Template, error) {
	templatePath := filepath.Join(l.templatesDir, name)

	// Check if template directory exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template not found: %s", name)
	}

	// Load template.yaml
	metadataPath := filepath.Join(templatePath, "template.yaml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template.yaml: %w", err)
	}

	// Parse YAML
	var tmpl Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse template.yaml: %w", err)
	}

	// Store template path
	tmpl.Path = templatePath

	// Validate template
	if err := l.validate(&tmpl); err != nil {
		return nil, fmt.Errorf("invalid template: %w", err)
	}

	return &tmpl, nil
}

// List returns all available templates
func (l *Loader) List() ([]string, error) {
	var templates []string

	// Walk through templates directory
	err := filepath.Walk(l.templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if this is a template.yaml file
		if !info.IsDir() && info.Name() == "template.yaml" {
			// Get relative path from templates dir
			relPath, err := filepath.Rel(l.templatesDir, filepath.Dir(path))
			if err != nil {
				return err
			}

			templates = append(templates, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	return templates, nil
}

// validate performs basic validation on a template
func (l *Loader) validate(tmpl *Template) error {
	if tmpl.Version == "" {
		return fmt.Errorf("version is required")
	}

	if tmpl.Name == "" {
		return fmt.Errorf("name is required")
	}

	if tmpl.Language == "" {
		return fmt.Errorf("language is required")
	}

	// Validate that all file sources exist
	filesDir := filepath.Join(tmpl.Path, "files")
	for _, file := range tmpl.Files {
		filePath := filepath.Join(filesDir, file.Source)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", file.Source)
		}
	}

	return nil
}

// GetFilesDir returns the files directory for a template
func (l *Loader) GetFilesDir(tmpl *Template) string {
	return filepath.Join(tmpl.Path, "files")
}
