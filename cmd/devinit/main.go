package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/renan-dev/devinit/internal/generator"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "devinit",
		Short: "Multi-language project scaffolding CLI",
		Long: `devinit is a CLI tool that creates production-ready projects
for multiple languages and frameworks with standardized structure,
Docker support, and best practices built-in.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Add subcommands
	rootCmd.AddCommand(newNewCmd())
	rootCmd.AddCommand(newValidateCmd())
	rootCmd.AddCommand(newDoctorCmd())
	rootCmd.AddCommand(newTemplatesCmd())

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")

	return rootCmd
}

func newNewCmd() *cobra.Command {
	var (
		lang        string
		framework   string
		docker      bool
		database    string
		ci          string
		noValidate  bool
		dryRun      bool
		pythonVersion string
		includeTests  bool
	)

	cmd := &cobra.Command{
		Use:   "new [type] [name]",
		Short: "Create a new project",
		Long: `Create a new project with the specified language and framework.

Examples:
  # Interactive mode
  devinit new

  # Non-interactive mode
  devinit new api my-service --lang python --framework fastapi

  # With all options
  devinit new api my-service \
    --lang python \
    --framework fastapi \
    --docker \
    --database postgres \
    --ci github`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNewCommand(args, lang, framework, database, pythonVersion, docker, includeTests, dryRun)
		},
	}

	cmd.Flags().StringVar(&lang, "lang", "", "programming language (python, nodejs, kotlin)")
	cmd.Flags().StringVar(&framework, "framework", "", "framework to use")
	cmd.Flags().BoolVar(&docker, "docker", true, "include Docker configuration")
	cmd.Flags().StringVar(&database, "database", "none", "database to configure (postgres, sqlite, none)")
	cmd.Flags().StringVar(&ci, "ci", "", "CI provider (github, gitlab, none)")
	cmd.Flags().BoolVar(&noValidate, "no-validate", false, "skip validation")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be done without doing it")
	cmd.Flags().StringVar(&pythonVersion, "python-version", "3.11", "Python version (python only)")
	cmd.Flags().BoolVar(&includeTests, "tests", true, "include test setup")

	return cmd
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate project structure",
		Long:  "Validate that the current project follows devinit standards",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement validation
			fmt.Println("Validating project...")
			return nil
		},
	}
}

func newDoctorCmd() *cobra.Command {
	var templateName string

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check system requirements",
		Long:  "Check that all required system dependencies are installed",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement doctor checks
			fmt.Println("Checking system requirements...")
			return nil
		},
	}

	cmd.Flags().StringVar(&templateName, "template", "", "check requirements for specific template")

	return cmd
}

func newTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage templates",
		Long:  "List, show, and manage project templates",
	}

	cmd.AddCommand(newTemplatesListCmd())
	cmd.AddCommand(newTemplatesShowCmd())
	cmd.AddCommand(newTemplatesValidateCmd())

	return cmd
}

func newTemplatesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			gen := getGenerator()
			templates, err := gen.ListTemplates()
			if err != nil {
				return err
			}

			fmt.Println("Available templates:")
			for _, tmpl := range templates {
				fmt.Printf("  - %s\n", tmpl)
			}
			return nil
		},
	}
}

func newTemplatesShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show [template]",
		Short: "Show template details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gen := getGenerator()
			tmpl, err := gen.GetTemplate(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Name: %s\n", tmpl.Name)
			fmt.Printf("Version: %s\n", tmpl.Version)
			fmt.Printf("Description: %s\n", tmpl.Description)
			fmt.Printf("Language: %s\n", tmpl.Language)
			fmt.Printf("Framework: %s\n", tmpl.Framework)
			fmt.Println("\nVariables:")
			for key, variable := range tmpl.Variables {
				fmt.Printf("  %s (%s): %s\n", key, variable.Type, variable.Description)
			}
			return nil
		},
	}
}

func newTemplatesValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate all templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			gen := getGenerator()
			templates, err := gen.ListTemplates()
			if err != nil {
				return err
			}

			fmt.Println("Validating templates...")
			errors := 0
			for _, name := range templates {
				_, err := gen.GetTemplate(name)
				if err != nil {
					fmt.Printf("  ✗ %s: %v\n", name, err)
					errors++
				} else {
					fmt.Printf("  ✓ %s\n", name)
				}
			}

			if errors > 0 {
				return fmt.Errorf("%d template(s) failed validation", errors)
			}

			fmt.Println("\nAll templates valid!")
			return nil
		},
	}
}

// Helper functions

func getTemplatesDir() string {
	// Get executable directory
	exe, err := os.Executable()
	if err != nil {
		// Fallback to current directory
		return "templates"
	}

	exeDir := filepath.Dir(exe)

	// Check if templates directory exists relative to executable
	templatesDir := filepath.Join(exeDir, "..", "templates")
	if _, err := os.Stat(templatesDir); err == nil {
		return templatesDir
	}

	// Fallback to templates in current directory (development mode)
	return "templates"
}

func getGenerator() *generator.Generator {
	return generator.NewGenerator(getTemplatesDir())
}

func runNewCommand(args []string, lang, framework, database, pythonVersion string, docker, includeTests, dryRun bool) error {
	// Determine project name
	projectName := ""
	if len(args) >= 2 {
		projectName = args[1]
	} else if len(args) == 1 {
		projectName = args[0]
	} else {
		return fmt.Errorf("project name is required")
	}

	// Determine language and framework
	if lang == "" {
		return fmt.Errorf("--lang flag is required")
	}

	if framework == "" {
		return fmt.Errorf("--framework flag is required")
	}

	// Build variables
	variables := map[string]interface{}{
		"ProjectName":    projectName,
		"PythonVersion":  pythonVersion,
		"IncludeDocker":  docker,
		"Database":       database,
		"IncludeTests":   includeTests,
	}

	// Create generator options
	opts := &generator.Options{
		ProjectName: projectName,
		Language:    lang,
		Framework:   framework,
		Variables:   variables,
		DryRun:      dryRun,
	}

	// Generate project
	gen := getGenerator()

	fmt.Printf("Creating %s/%s project: %s\n", lang, framework, projectName)
	if dryRun {
		fmt.Println("(dry run - no files will be created)")
	}

	if err := gen.Generate(opts); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	if !dryRun {
		fmt.Printf("\n✓ Project created successfully at: ./%s\n", projectName)
		fmt.Println("\nNext steps:")
		fmt.Printf("  cd %s\n", projectName)

		if lang == "python" {
			fmt.Println("  poetry install")
			if docker {
				fmt.Println("  docker compose up")
			} else {
				fmt.Println("  poetry run uvicorn src.main:app --reload")
			}
		}
	}

	return nil
}
