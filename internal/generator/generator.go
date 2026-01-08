package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/renan-dev/devinit/internal/template"
)

// Generator generates projects from templates
type Generator struct {
	loader   *template.Loader
	renderer *template.Renderer
}

// NewGenerator creates a new project generator
func NewGenerator(templatesDir string) *Generator {
	return &Generator{
		loader:   template.NewLoader(templatesDir),
		renderer: template.NewRenderer(),
	}
}

// Options for project generation
type Options struct {
	ProjectName string
	Language    string
	Framework   string
	OutputDir   string
	Variables   map[string]interface{}
	DryRun      bool
}

// Generate creates a new project from a template
func (g *Generator) Generate(opts *Options) error {
	// Construct template name
	templateName := fmt.Sprintf("%s/%s", opts.Language, opts.Framework)

	// Load template
	tmpl, err := g.loader.Load(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Merge options with template variables
	variables := g.mergeVariables(tmpl, opts.Variables)

	// Create context
	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = opts.ProjectName
	}

	ctx := template.NewContext(opts.ProjectName, outputDir, variables, tmpl)

	// Create project directory
	if !opts.DryRun {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
	}

	// Generate files
	filesDir := g.loader.GetFilesDir(tmpl)
	for _, fileSpec := range tmpl.Files {
		// Check if file should be generated based on conditions
		if !g.shouldGenerateFile(fileSpec, ctx) {
			if opts.DryRun {
				fmt.Printf("Skipped: %s (conditions not met)\n", fileSpec.Destination)
			}
			continue
		}

		if err := g.generateFile(filesDir, fileSpec, ctx, opts.DryRun); err != nil {
			return fmt.Errorf("failed to generate file %s: %w", fileSpec.Destination, err)
		}
	}

	if !opts.DryRun {
		// Create .devinit.yaml metadata file
		if err := g.createMetadataFile(ctx, tmpl); err != nil {
			return fmt.Errorf("failed to create metadata file: %w", err)
		}
	}

	return nil
}

// generateFile generates a single file from template
func (g *Generator) generateFile(filesDir string, fileSpec template.FileSpec, ctx *template.Context, dryRun bool) error {
	sourcePath := filepath.Join(filesDir, fileSpec.Source)
	destPath := filepath.Join(ctx.OutputDir, fileSpec.Destination)

	// Check if file should be rendered
	if g.renderer.ShouldRender(fileSpec.Source) {
		// Get actual output filename (without .tmpl)
		actualDest := filepath.Join(ctx.OutputDir, g.renderer.GetOutputFilename(fileSpec.Destination))

		if dryRun {
			fmt.Printf("Would render: %s -> %s\n", fileSpec.Source, actualDest)
			return nil
		}

		// Render template
		if err := g.renderer.RenderToFile(sourcePath, actualDest, ctx, fileSpec.GetPermissions()); err != nil {
			return err
		}

		fmt.Printf("Created: %s\n", actualDest)
	} else {
		if dryRun {
			fmt.Printf("Would copy: %s -> %s\n", fileSpec.Source, destPath)
			return nil
		}

		// Copy static file
		if err := g.renderer.CopyFile(sourcePath, destPath, fileSpec.GetPermissions()); err != nil {
			return err
		}

		fmt.Printf("Created: %s\n", destPath)
	}

	return nil
}

// shouldGenerateFile checks if a file should be generated based on its conditions
func (g *Generator) shouldGenerateFile(fileSpec template.FileSpec, ctx *template.Context) bool {
	if len(fileSpec.Conditions) == 0 {
		return true
	}

	for _, condition := range fileSpec.Conditions {
		if !g.evaluateCondition(condition, ctx) {
			return false
		}
	}

	return true
}

// evaluateCondition evaluates a single condition string
// Supports: {{ .VariableName }}, variable names, and simple expressions
func (g *Generator) evaluateCondition(condition string, ctx *template.Context) bool {
	condition = strings.TrimSpace(condition)

	condition = strings.TrimSpace(condition)
	if strings.HasPrefix(condition, "{{") && strings.HasSuffix(condition, "}}") {
		condition = strings.TrimSpace(condition[2 : len(condition)-2])
	}

	condition = strings.TrimPrefix(condition, ".")

	switch condition {
	case "IncludeDocker":
		return ctx.IncludeDocker
	case "IncludeTests":
		return ctx.IncludeTests
	}

	return ctx.GetBool(condition)
}

// mergeVariables merges user-provided variables with template defaults
func (g *Generator) mergeVariables(tmpl *template.Template, userVars map[string]interface{}) map[string]interface{} {
	variables := make(map[string]interface{})

	// Start with template defaults
	for key, varDef := range tmpl.Variables {
		if varDef.Default != nil {
			variables[key] = varDef.Default
		}
	}

	// Override with user-provided values
	for key, value := range userVars {
		variables[key] = value
	}

	return variables
}

// createMetadataFile creates the .devinit.yaml file in the project
func (g *Generator) createMetadataFile(ctx *template.Context, tmpl *template.Template) error {
	metadata := fmt.Sprintf(`schema_version: "1.0"
template:
  name: %s/%s
  version: %s
variables:
`, tmpl.Language, tmpl.Framework, tmpl.Version)

	for key, value := range ctx.Variables {
		metadata += fmt.Sprintf("  %s: %v\n", key, value)
	}

	metadataPath := filepath.Join(ctx.OutputDir, ".devinit.yaml")
	return os.WriteFile(metadataPath, []byte(metadata), 0644)
}

// ListTemplates returns all available templates
func (g *Generator) ListTemplates() ([]string, error) {
	return g.loader.List()
}

// GetTemplate returns a specific template
func (g *Generator) GetTemplate(name string) (*template.Template, error) {
	return g.loader.Load(name)
}
