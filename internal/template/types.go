package template

import "os"

// Template represents a project template
type Template struct {
	// Metadata
	Version     string `yaml:"version"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Language    string `yaml:"language"`
	Framework   string `yaml:"framework"`
	MinCLIVersion string `yaml:"min_cli_version"`

	// Requirements
	Requirements Requirements `yaml:"requirements"`

	// Variables
	Variables map[string]Variable `yaml:"variables"`

	// Files
	Files []FileSpec `yaml:"files"`

	// Dependencies on other templates
	Dependencies []Dependency `yaml:"dependencies"`

	// Lifecycle hooks
	Hooks Hooks `yaml:"hooks"`

	// Healthcheck configuration
	Healthcheck *Healthcheck `yaml:"healthcheck,omitempty"`

	// Internal fields (not in YAML)
	Path string `yaml:"-"` // Path to template directory
}

// Requirements defines system requirements
type Requirements struct {
	System      []SystemRequirement      `yaml:"system,omitempty"`
	Environment []EnvironmentRequirement `yaml:"environment,omitempty"`
}

// SystemRequirement represents a required command/binary
type SystemRequirement struct {
	Command     string `yaml:"command"`
	Version     string `yaml:"version,omitempty"`
	Required    bool   `yaml:"required"`
	When        string `yaml:"when,omitempty"`
	InstallHint string `yaml:"install_hint,omitempty"`
}

// EnvironmentRequirement represents required environment variable
type EnvironmentRequirement struct {
	Variable string `yaml:"var"`
	Required bool   `yaml:"required"`
	When     string `yaml:"when,omitempty"`
}

// VariableType represents the type of a template variable
type VariableType string

const (
	VariableTypeString  VariableType = "string"
	VariableTypeBool    VariableType = "boolean"
	VariableTypeChoice  VariableType = "choice"
	VariableTypeInt     VariableType = "int"
)

// Variable defines a template variable
type Variable struct {
	Type        VariableType `yaml:"type"`
	Required    bool         `yaml:"required"`
	Default     interface{}  `yaml:"default,omitempty"`
	Choices     []string     `yaml:"choices,omitempty"`
	Pattern     string       `yaml:"pattern,omitempty"`
	Description string       `yaml:"description,omitempty"`
}

// FileSpec specifies a file to be generated
type FileSpec struct {
	Source      string   `yaml:"src"`
	Destination string   `yaml:"dest"`
	Conditions  []string `yaml:"conditions,omitempty"`
	Permissions string   `yaml:"permissions,omitempty"`
}

// GetPermissions returns the file permissions as os.FileMode
func (f *FileSpec) GetPermissions() os.FileMode {
	if f.Permissions == "" {
		return 0644 // default
	}
	// TODO: parse octal string to os.FileMode
	return 0644
}

// Dependency represents a dependency on another template
type Dependency struct {
	Template string `yaml:"template"`
	When     string `yaml:"when,omitempty"`
}

// Hooks defines lifecycle hooks
type Hooks struct {
	PreGenerate  []Hook `yaml:"pre_generate,omitempty"`
	PostGenerate []Hook `yaml:"post_generate,omitempty"`
}

// ErrorLevel represents how to handle hook errors
type ErrorLevel string

const (
	ErrorLevelError  ErrorLevel = "error"
	ErrorLevelWarn   ErrorLevel = "warn"
	ErrorLevelIgnore ErrorLevel = "ignore"
)

// Hook represents a lifecycle hook command
type Hook struct {
	Run        string     `yaml:"run,omitempty"`
	Validate   string     `yaml:"validate,omitempty"`
	WorkingDir string     `yaml:"working_dir,omitempty"`
	ErrorLevel ErrorLevel `yaml:"error_level,omitempty"`
	Error      string     `yaml:"error,omitempty"` // Custom error message
}

// Healthcheck defines healthcheck configuration for generated project
type Healthcheck struct {
	Command string `yaml:"command"`
	Port    int    `yaml:"port"`
	Timeout string `yaml:"timeout,omitempty"`
}

// Context represents the context for template rendering
type Context struct {
	// Project information
	ProjectName string
	OutputDir   string

	// Template variables (as map for programmatic access)
	Variables map[string]interface{}

	// Template reference
	Template *Template

	// Computed values
	ProjectNameSnake  string
	ProjectNameCamel  string
	ProjectNamePascal string
	ProjectNameKebab  string

	// Common template variables (exposed as fields for easy template access)
	PythonVersion  string
	IncludeDocker  bool
	Database       string
	IncludeTests   bool
	CIProvider     string
}

// NewContext creates a new template context
func NewContext(projectName, outputDir string, variables map[string]interface{}, tmpl *Template) *Context {
	ctx := &Context{
		ProjectName:       projectName,
		OutputDir:         outputDir,
		Variables:         variables,
		Template:          tmpl,
		ProjectNameSnake:  toSnakeCase(projectName),
		ProjectNameCamel:  toCamelCase(projectName),
		ProjectNamePascal: toPascalCase(projectName),
		ProjectNameKebab:  toKebabCase(projectName),
	}

	// Extract common variables to fields for template access
	if v, ok := variables["PythonVersion"].(string); ok {
		ctx.PythonVersion = v
	}
	if v, ok := variables["IncludeDocker"].(bool); ok {
		ctx.IncludeDocker = v
	}
	if v, ok := variables["Database"].(string); ok {
		ctx.Database = v
	}
	if v, ok := variables["IncludeTests"].(bool); ok {
		ctx.IncludeTests = v
	}
	if v, ok := variables["CIProvider"].(string); ok {
		ctx.CIProvider = v
	}

	return ctx
}

// GetString retrieves a string variable value
func (c *Context) GetString(key string) string {
	if v, ok := c.Variables[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetBool retrieves a boolean variable value
func (c *Context) GetBool(key string) bool {
	if v, ok := c.Variables[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// GetInt retrieves an integer variable value
func (c *Context) GetInt(key string) int {
	if v, ok := c.Variables[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}
