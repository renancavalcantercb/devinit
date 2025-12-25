package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Renderer renders template files
type Renderer struct {
	funcMap template.FuncMap
}

// NewRenderer creates a new template renderer
func NewRenderer() *Renderer {
	funcMap := template.FuncMap{
		// String manipulation
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
		"title":   strings.Title,
		"snake":   toSnakeCase,
		"camel":   toCamelCase,
		"pascal":  toPascalCase,
		"kebab":   toKebabCase,

		// String operations
		"contains": strings.Contains,
		"replace":  strings.ReplaceAll,
		"trim":     strings.TrimSpace,
		"split":    strings.Split,
		"join":     strings.Join,

		// Comparison
		"eq": func(a, b interface{}) bool { return a == b },
		"ne": func(a, b interface{}) bool { return a != b },
	}

	return &Renderer{
		funcMap: funcMap,
	}
}

// Render renders a single template file
func (r *Renderer) Render(templatePath string, ctx *Context) (string, error) {
	// Read template content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Create template
	tmpl, err := template.New(filepath.Base(templatePath)).
		Funcs(r.funcMap).
		Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderToFile renders a template and writes it to a file
func (r *Renderer) RenderToFile(templatePath, outputPath string, ctx *Context, perm os.FileMode) error {
	// Render template
	content, err := r.Render(templatePath, ctx)
	if err != nil {
		return err
	}

	// Create parent directory if needed
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(outputPath, []byte(content), perm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// CopyFile copies a static file (no template rendering)
func (r *Renderer) CopyFile(srcPath, dstPath string, perm os.FileMode) error {
	// Read source
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create parent directory if needed
	dir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to destination
	if err := os.WriteFile(dstPath, content, perm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ShouldRender returns true if the file should be rendered (has .tmpl extension)
func (r *Renderer) ShouldRender(filename string) bool {
	return strings.HasSuffix(filename, ".tmpl")
}

// GetOutputFilename returns the output filename (removes .tmpl extension if present)
func (r *Renderer) GetOutputFilename(filename string) string {
	if strings.HasSuffix(filename, ".tmpl") {
		return strings.TrimSuffix(filename, ".tmpl")
	}
	return filename
}
