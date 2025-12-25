# devinit - Architecture Decision Record

**Version**: 1.0.0
**Status**: Approved
**Last Updated**: 2025-12-25
**Authors**: Development Team

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Product Vision](#product-vision)
3. [Architecture Decisions](#architecture-decisions)
4. [Technical Specifications](#technical-specifications)
5. [Implementation Details](#implementation-details)
6. [Quality Assurance](#quality-assurance)
7. [Roadmap](#roadmap)

---

## Executive Summary

devinit is a multi-language project scaffolding CLI tool designed to reduce project initialization time from hours to seconds while enforcing production-ready standards across different technology stacks.

**Key Metrics**:
- Project creation: < 60 seconds
- Zero manual edits to run locally
- Startup time: < 200ms
- Binary size target: < 10MB

**Core Principles**:
- Opinionated templates for consistency
- Docker-first approach
- Production-ready from day one
- Multi-language support with unified interface

---

## Product Vision

### Target Personas

**Backend Developer**
- Creates new services frequently
- Values speed without sacrificing quality
- Expects Docker, clean structure, clear README

**DevOps / Platform Engineer**
- Enforces standards across projects
- Reduces structural variation
- Facilitates CI/CD and observability

### Primary User Story

```
As a developer or DevOps engineer
I want to create new projects through a simple, interactive CLI
So that I can start services in different languages with standardized,
production-ready structure
```

### Success Criteria

1. Create functional project in < 60 seconds
2. Zero manual edits to run locally
3. Adoption in real projects (internal or OSS)
4. Template reusability across languages

---

## Architecture Decisions

### ADR-001: Implementation Language

**Decision**: Go

**Context**: CLI requires fast execution, easy distribution, and maintainability.

**Rationale**:
- Single binary distribution (no runtime dependencies)
- Excellent CLI ecosystem (cobra, viper, bubble tea)
- Fast startup time (~100ms)
- Cross-compilation support (Linux, macOS, Windows)
- Balance between development speed and performance

**Alternatives Considered**:
- Rust: Better performance but slower development, steeper learning curve
- Python: Faster development but complex distribution, slower startup

**Consequences**:
- Team needs Go knowledge
- Excellent tooling support
- Easy CI/CD integration
- Native cross-platform builds

---

### ADR-002: Project Structure

**Decision**: Monorepo with embedded templates

**Structure**:
```
devinit/
├── cmd/
│   └── devinit/
│       └── main.go                 # CLI entrypoint
├── internal/
│   ├── generator/
│   │   ├── generator.go            # Core generator interface
│   │   ├── python.go               # Python-specific logic
│   │   ├── nodejs.go               # Node.js-specific logic
│   │   └── kotlin.go               # Kotlin-specific logic
│   ├── template/
│   │   ├── loader.go               # Template loading
│   │   ├── renderer.go             # Go template rendering
│   │   ├── validator.go            # Template validation
│   │   └── metadata.go             # Template metadata parsing
│   ├── config/
│   │   ├── config.go               # Global config (~/.config/devinit/)
│   │   └── project.go              # Project metadata (.devinit.yaml)
│   ├── prompt/
│   │   ├── interactive.go          # Interactive wizard
│   │   └── validators.go           # Input validation
│   ├── registry/
│   │   └── registry.go             # Template registry (future)
│   └── validator/
│       ├── system.go               # System requirements validation
│       └── project.go              # Generated project validation
├── templates/
│   ├── python/
│   │   ├── fastapi/
│   │   │   ├── template.yaml       # Template metadata
│   │   │   └── files/              # Template files
│   │   │       ├── main.py.tmpl
│   │   │       ├── Dockerfile
│   │   │       ├── README.md.tmpl
│   │   │       ├── .env.example
│   │   │       ├── .gitignore
│   │   │       └── pyproject.toml.tmpl
│   │   └── flask/
│   ├── nodejs/
│   │   └── express/
│   ├── kotlin/
│   │   └── spring-boot/
│   └── shared/                      # Reusable components
│       ├── docker/
│       │   ├── Dockerfile.base
│       │   └── .dockerignore
│       ├── github-actions/
│       │   └── ci.yml.tmpl
│       └── prometheus/
│           └── prometheus.yml.tmpl
├── pkg/
│   └── devinit/                     # Public API (if needed)
├── test/
│   ├── unit/
│   ├── integration/
│   └── e2e/
│       ├── testdata/
│       └── generated/               # Generated projects for testing
├── docs/
│   ├── architecture/
│   ├── templates/
│   └── user-guide/
├── scripts/
│   ├── build.sh
│   └── test-all-templates.sh
├── .goreleaser.yml                  # Release automation
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

**Rationale**:
- Atomic releases (CLI + templates versioned together)
- Simplified CI/CD (single build, single test)
- No synchronization issues between repos
- Embedded templates via go:embed (no external dependencies)

**Migration Path**:
- When templates > 50MB or community contributions grow
- Split to separate repo with version pinning

---

### ADR-003: Template System

**Decision**: Hybrid approach (YAML metadata + Go templates)

**Template Structure**:
```
templates/python/fastapi/
├── template.yaml          # Metadata and configuration
└── files/                 # Template files
    ├── main.py.tmpl      # Go template syntax
    ├── Dockerfile        # Static file
    └── README.md.tmpl    # Go template syntax
```

**Template Metadata Schema**:
```yaml
version: "1.0.0"
name: "FastAPI API"
description: "Production-ready FastAPI service"

language: python
framework: fastapi
min_cli_version: "1.0.0"

# Language-specific requirements
requirements:
  system:
    - command: python3
      version: ">=3.11"
      required: true
      install_hint: "https://www.python.org/downloads/"

    - command: docker
      version: ">=24.0"
      required: false
      when: "{{ .IncludeDocker }}"
      install_hint: "https://docs.docker.com/install/"

    - command: poetry
      required: true
      install_hint: "curl -sSL https://install.python-poetry.org | python3 -"

# User-facing variables
variables:
  project_name:
    type: string
    required: true
    pattern: "^[a-z][a-z0-9-]*$"
    description: "Project name (lowercase, hyphens allowed)"

  python_version:
    type: string
    default: "3.11"
    choices: ["3.11", "3.12", "3.13"]
    description: "Python version"

  include_docker:
    type: boolean
    default: true
    description: "Include Docker configuration"

  database:
    type: choice
    choices: ["postgres", "mysql", "sqlite", "none"]
    default: "postgres"
    description: "Database to configure"

  include_tests:
    type: boolean
    default: true
    description: "Include pytest setup"

# File generation rules
files:
  - src: "main.py.tmpl"
    dest: "src/main.py"

  - src: "Dockerfile"
    dest: "Dockerfile"
    conditions: ["include_docker"]

  - src: "docker-compose.yml.tmpl"
    dest: "docker-compose.yml"
    conditions: ["include_docker", "database != 'none'"]

  - src: "tests/test_main.py.tmpl"
    dest: "tests/test_main.py"
    conditions: ["include_tests"]

# Shared component dependencies
dependencies:
  - template: "shared/docker"
    when: "{{ .IncludeDocker }}"

  - template: "shared/github-actions"
    when: "{{ .CIProvider == 'github' }}"

# Lifecycle hooks
hooks:
  pre_generate:
    - validate: "python3 --version"
      error: "Python 3 not found"

  post_generate:
    - run: "poetry install"
      working_dir: "{{ .OutputDir }}"
      error_level: "warn"

    - run: "poetry run black ."
      working_dir: "{{ .OutputDir }}"
      error_level: "warn"

    - run: "git init"
      working_dir: "{{ .OutputDir }}"
      error_level: "ignore"

# Health check configuration
healthcheck:
  command: "curl -f http://localhost:8000/health"
  port: 8000
  timeout: 5s

# Documentation
documentation:
  readme_sections:
    - setup
    - development
    - deployment
    - testing
```

**Go Template Example**:
```python
# templates/python/fastapi/files/main.py.tmpl
from fastapi import FastAPI
{{- if .Database }}
from sqlalchemy import create_engine
{{- end }}

app = FastAPI(
    title="{{ .ProjectName }}",
    version="0.1.0"
)

{{- if .IncludeHealthcheck }}

@app.get("/health")
async def health():
    return {"status": "healthy"}
{{- end }}

{{- if .Database }}

# Database configuration
DATABASE_URL = "{{ .DatabaseURL }}"
engine = create_engine(DATABASE_URL)
{{- end }}
```

**Rationale**:
- Clear separation between metadata and content
- Go templates: powerful, well-tested, no external dependencies
- YAML metadata: human-readable, versionable, schema-validatable
- Conditional rendering without code compilation

---

### ADR-004: CLI Interface Design

**Decision**: Cobra-based CLI with interactive and non-interactive modes

**Command Structure**:
```bash
devinit [command] [flags]

Commands:
  new         Create a new project
  add         Add components to existing project
  validate    Validate project structure
  upgrade     Upgrade project to newer template version
  templates   Manage templates
  doctor      Check system requirements
  version     Show version information

Flags:
  -h, --help      Help for command
  -v, --verbose   Verbose output
  --no-validate   Skip validation
  --dry-run       Show what would be done without doing it
```

**Primary Commands**:

```bash
# Interactive mode (default)
devinit new

# Non-interactive (CI/scripts)
devinit new api my-service \
  --lang python \
  --framework fastapi \
  --docker \
  --database postgres \
  --ci github

# Add components
devinit add docker
devinit add ci --provider github
devinit add monitoring --stack prometheus

# Template management
devinit templates list
devinit templates show python/fastapi
devinit templates validate
devinit templates add https://github.com/company/templates  # future

# Project operations
devinit validate                    # Validate current project
devinit upgrade --dry-run          # Check available upgrades
devinit upgrade --interactive      # Interactive upgrade

# Diagnostics
devinit doctor                     # Check all requirements
devinit doctor --template python/fastapi  # Check specific template
devinit version
```

**Configuration Files**:

Global config (`~/.config/devinit/config.yaml`):
```yaml
version: "1.0"

# Template sources (future feature)
template_registry:
  - url: https://github.com/devinit/templates
    branch: main
  - url: https://github.com/mycompany/templates
    token: ${GITHUB_TOKEN}

# Default values for new projects
defaults:
  language: python
  docker: true
  ci_provider: github
  git_init: true

  python:
    package_manager: poetry
    min_version: "3.11"
    formatter: black

  nodejs:
    package_manager: pnpm
    min_version: "20"
    formatter: prettier

# Project metadata defaults
project_defaults:
  author: "Your Name"
  email: "you@example.com"
  license: MIT
  organization: ""

# Behavior
behavior:
  interactive: true           # Default to interactive mode
  validation_level: basic     # basic, strict, none
  auto_git_init: true
  auto_install_deps: ask      # yes, no, ask
```

Project metadata (`.devinit.yaml` in generated project):
```yaml
schema_version: "1.0"

# Template information
template:
  name: python/fastapi
  version: "1.2.0"
  cli_version: "1.0.0"

# Generation timestamp
generated_at: "2025-12-25T10:00:00Z"
updated_at: "2025-12-25T10:00:00Z"

# Variables used during generation
variables:
  project_name: my-service
  python_version: "3.11"
  include_docker: true
  database: postgres
  include_tests: true

# User customizations (tracked for upgrades)
customizations:
  added:
    - docker-compose.override.yml
    - src/custom_module.py

  modified:
    - pyproject.toml
    - README.md

  ignored:
    - src/main.py  # Don't upgrade this file

# Upgrade history
upgrades:
  - from: "1.0.0"
    to: "1.2.0"
    date: "2025-12-26T15:30:00Z"
    type: minor
```

---

### ADR-005: Versioning Strategy

**Decision**: Semantic versioning with separate CLI and template versions

**Version Scheme**:
```
CLI Version:      v1.2.3
Template Version: v1.2.3

Major version must match:
  CLI v1.x.x  ↔  Templates v1.x.x  ✓
  CLI v1.x.x  ↔  Templates v2.x.x  ✗
```

**Breaking Change Policy**:
```
Template v1.2.3
         │ │ │
         │ │ └─ PATCH: Bug fixes, typos, documentation
         │ │           - No structural changes
         │ │           - Safe to auto-upgrade
         │ │
         │ └─── MINOR: New features, optional additions
         │             - New optional files
         │             - New optional variables
         │             - Backwards compatible
         │             - Safe to upgrade with review
         │
         └───── MAJOR: Breaking changes
                       - Changed directory structure
                       - Removed variables
                       - Changed file locations
                       - Requires manual migration
```

**Examples**:
```
v1.0.0 → v1.0.1: Fixed typo in README                    (patch)
v1.0.1 → v1.1.0: Added optional .dockerignore            (minor)
v1.1.0 → v1.2.0: Added healthcheck endpoint              (minor)
v1.2.0 → v2.0.0: Changed src/ to app/ directory          (major)
```

**Compatibility Matrix**:
```yaml
# In template.yaml
min_cli_version: "1.0.0"    # Minimum required CLI version
max_cli_version: "2.0.0"    # Optional: for deprecation warnings
```

**Upgrade Strategy**:
```bash
# Check for upgrades
devinit upgrade --check
# Output:
# Current: python/fastapi v1.0.0
# Available: v1.2.0 (minor - safe)
#            v2.0.0 (major - review required)

# Dry run
devinit upgrade --dry-run --to v1.2.0
# Shows: files to be added, modified, conflicts

# Interactive upgrade
devinit upgrade --interactive
# Walks through changes, asks for confirmation

# Automatic patch/minor upgrades
devinit upgrade --auto-minor
```

---

### ADR-006: Validation System

**Decision**: Three-tier validation (pre-flight, runtime, post-generation)

**Validation Levels**:
```go
type ValidationLevel int

const (
    ValidationNone   ValidationLevel = iota  // --no-validate
    ValidationBasic                          // default
    ValidationStrict                         // --strict (CI mode)
)
```

**Tier 1: Pre-flight Checks**
```
Runs before project generation
Checks system requirements
Non-blocking warnings, blocking errors
```

Example output:
```
Checking system requirements...
  ✓ git found (v2.43.0)
  ✗ docker not found
    → Install: https://docs.docker.com/install
  ⚠ python found but v3.9 (recommended: v3.11+)
    → Upgrade: https://www.python.org/downloads

1 error, 1 warning found.
Continue anyway? [y/N]
```

**Tier 2: Runtime Validation**
```
Runs during template rendering
Validates template variables
Checks conditional logic
```

**Tier 3: Post-generation Validation**
```
Runs after project created
Validates generated structure
Tests basic functionality
```

Example:
```bash
devinit new api my-service --lang python

# Auto-runs post-generation validation
✓ Generated 15 files
✓ Running validation...
  ✓ Directory structure correct
  ✓ Dockerfile syntax valid
  ✓ Python imports resolve
  ✓ pyproject.toml valid
  ⚠ Poetry not installed - run: pip install poetry

Project created at: ./my-service

Next steps:
  cd my-service
  poetry install
  docker compose up
```

**Implementation**:
```go
// internal/validator/system.go
type SystemValidator struct {
    Level ValidationLevel
}

type Requirement struct {
    Command     string
    Version     string  // semver constraint: >=3.11, ^20.0.0
    Required    bool    // block if not found
    When        string  // condition: "{{ .IncludeDocker }}"
    InstallHint string
}

func (v *SystemValidator) Validate(reqs []Requirement) error {
    var errors []error
    var warnings []error

    for _, req := range reqs {
        if !req.ShouldCheck() {
            continue
        }

        exists, version := checkCommand(req.Command)

        if !exists {
            err := fmt.Errorf("%s not found\n  %s",
                req.Command, req.InstallHint)

            if req.Required {
                errors = append(errors, err)
            } else {
                warnings = append(warnings, err)
            }
            continue
        }

        if req.Version != "" && !versionMatches(version, req.Version) {
            err := fmt.Errorf("%s version %s does not match %s",
                req.Command, version, req.Version)

            if v.Level == ValidationStrict {
                errors = append(errors, err)
            } else {
                warnings = append(warnings, err)
            }
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("validation failed:\n%v", errors)
    }

    return nil
}
```

---

### ADR-007: Testing Strategy

**Decision**: Multi-layered testing (unit, integration, E2E)

**Test Structure**:
```
test/
├── unit/                           # Fast, isolated tests
│   ├── generator_test.go
│   ├── template_test.go
│   ├── validator_test.go
│   └── config_test.go
│
├── integration/                    # Component interaction
│   ├── cli_test.go
│   ├── template_rendering_test.go
│   └── validation_flow_test.go
│
└── e2e/                           # Full workflow tests
    ├── python_fastapi_test.go
    ├── nodejs_express_test.go
    ├── kotlin_springboot_test.go
    └── testdata/
        └── expected/              # Expected outputs
            ├── python-fastapi/
            └── nodejs-express/
```

**E2E Test Pattern**:
```go
func TestPythonFastAPI(t *testing.T) {
    tmpDir := t.TempDir()

    // 1. Generate project
    cmd := exec.Command("devinit", "new", "api", "test-service",
        "--lang", "python",
        "--framework", "fastapi",
        "--docker",
        "--no-interactive",
        "--no-validate")
    cmd.Dir = tmpDir
    output, err := cmd.CombinedOutput()
    require.NoError(t, err, "generation failed: %s", output)

    projectDir := filepath.Join(tmpDir, "test-service")

    // 2. Validate structure
    assertFileExists(t, projectDir, "src/main.py")
    assertFileExists(t, projectDir, "Dockerfile")
    assertFileExists(t, projectDir, "pyproject.toml")
    assertFileExists(t, projectDir, ".devinit.yaml")
    assertFileContains(t, projectDir, "src/main.py", "FastAPI")

    // 3. Validate metadata
    metadata := readDevinitYaml(t, projectDir)
    assert.Equal(t, "python/fastapi", metadata.Template.Name)
    assert.Equal(t, "test-service", metadata.Variables["project_name"])

    // 4. Test Docker build
    dockerBuild := exec.Command("docker", "build", "-t", "test-service", ".")
    dockerBuild.Dir = projectDir
    output, err = dockerBuild.CombinedOutput()
    require.NoError(t, err, "docker build failed: %s", output)

    // 5. Test container run + healthcheck
    dockerRun := exec.Command("docker", "run", "-d", "-p", "8000:8000",
        "--name", "test-service", "test-service")
    output, err = dockerRun.CombinedOutput()
    require.NoError(t, err, "docker run failed: %s", output)

    defer func() {
        exec.Command("docker", "rm", "-f", "test-service").Run()
    }()

    // Wait for service to start
    time.Sleep(2 * time.Second)

    // 6. Test healthcheck endpoint
    resp, err := http.Get("http://localhost:8000/health")
    require.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)

    body, _ := io.ReadAll(resp.Body)
    assert.Contains(t, string(body), "healthy")
}
```

**Template Validation Tests**:
```go
func TestAllTemplatesValid(t *testing.T) {
    templates := findAllTemplates("../../templates")

    for _, tmpl := range templates {
        t.Run(tmpl.Name, func(t *testing.T) {
            // Validate metadata schema
            err := validateTemplateYaml(tmpl.MetadataPath)
            assert.NoError(t, err)

            // Validate all template files render
            ctx := getTestContext()
            err = renderAllFiles(tmpl, ctx)
            assert.NoError(t, err)

            // Validate no broken references
            err = validateReferences(tmpl)
            assert.NoError(t, err)
        })
    }
}
```

**CI Pipeline**:
```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  unit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Run unit tests
        run: make test-unit

      - name: Coverage
        run: make coverage

  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5

      - name: Run integration tests
        run: make test-integration

  e2e:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        template:
          - python/fastapi
          - nodejs/express
          - kotlin/spring-boot

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5

      - name: Setup Docker
        uses: docker/setup-buildx-action@v3

      - name: Build CLI
        run: make build

      - name: Test ${{ matrix.template }}
        run: make test-e2e TEMPLATE=${{ matrix.template }}

  templates:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5

      - name: Validate all templates
        run: make validate-templates
```

---

### ADR-008: Extensibility (Post-MVP)

**Decision**: Plugin system for future extensibility (not in MVP)

**Future Plugin Interface**:
```go
// pkg/devinit/plugin.go
type Plugin interface {
    Name() string
    Version() string

    // Lifecycle hooks
    OnPreGenerate(ctx *Context) error
    OnPostGenerate(ctx *Context) error

    // Add custom commands
    Commands() []*cobra.Command

    // Add custom templates
    Templates() []*Template
}

// Example plugin
type TerraformPlugin struct{}

func (p *TerraformPlugin) Name() string { return "terraform" }

func (p *TerraformPlugin) OnPostGenerate(ctx *Context) error {
    // Generate Terraform files based on project
    return generateTerraform(ctx)
}

func (p *TerraformPlugin) Commands() []*cobra.Command {
    return []*cobra.Command{
        {
            Use:   "terraform",
            Short: "Generate Terraform configuration",
            Run:   terraformCommand,
        },
    }
}
```

**Plugin Loading** (future):
```bash
# Install plugin
devinit plugins install terraform

# Use plugin
devinit new api my-service --plugin terraform

# Plugin adds custom command
devinit terraform generate
```

**Template Registry** (future):
```yaml
# config.yaml
template_registry:
  - type: git
    url: https://github.com/devinit/templates
    branch: main

  - type: git
    url: git@github.com:company/private-templates
    auth:
      ssh_key: ~/.ssh/id_rsa

  - type: local
    path: ~/my-templates
```

---

## Technical Specifications

### Core Interfaces

```go
// internal/generator/generator.go
package generator

import "context"

type Generator interface {
    // Generate creates a new project
    Generate(ctx context.Context, opts *Options) error

    // Validate checks if generation is possible
    Validate(ctx context.Context, opts *Options) error

    // Language returns the language this generator handles
    Language() string
}

type Options struct {
    ProjectName  string
    Language     string
    Framework    string
    OutputDir    string
    Variables    map[string]interface{}
    Docker       bool
    CI           string
    Interactive  bool
    Validate     ValidationLevel
    DryRun       bool
}

type Context struct {
    ProjectName  string
    OutputDir    string
    Variables    map[string]interface{}
    Template     *Template
    Logger       Logger
}
```

```go
// internal/template/template.go
package template

type Template struct {
    // Metadata
    Version     string
    Name        string
    Description string
    Language    string
    Framework   string

    // Requirements
    Requirements Requirements

    // Variables
    Variables map[string]Variable

    // Files
    Files []FileSpec

    // Dependencies
    Dependencies []Dependency

    // Hooks
    Hooks Hooks

    // Healthcheck
    Healthcheck *Healthcheck
}

type Variable struct {
    Type        VariableType  // string, boolean, choice
    Required    bool
    Default     interface{}
    Choices     []string
    Pattern     string
    Description string
}

type FileSpec struct {
    Source      string
    Destination string
    Conditions  []string
    Permissions os.FileMode
}

type Hooks struct {
    PreGenerate  []Hook
    PostGenerate []Hook
}

type Hook struct {
    Run         string
    Validate    string
    WorkingDir  string
    ErrorLevel  ErrorLevel  // error, warn, ignore
}
```

```go
// internal/validator/validator.go
package validator

type Validator interface {
    Validate(ctx context.Context) error
}

type SystemValidator struct {
    Level        ValidationLevel
    Requirements []Requirement
}

type ProjectValidator struct {
    ProjectDir string
    Template   *Template
}

type Requirement struct {
    Command     string
    Version     string
    Required    bool
    When        string
    InstallHint string
}
```

### Template Rendering

```go
// internal/template/renderer.go
package template

import (
    "text/template"
)

type Renderer struct {
    template *template.Template
    funcMap  template.FuncMap
}

func NewRenderer() *Renderer {
    funcMap := template.FuncMap{
        // Custom functions
        "lower":     strings.ToLower,
        "upper":     strings.ToUpper,
        "snake":     toSnakeCase,
        "camel":     toCamelCase,
        "kebab":     toKebabCase,
        "contains":  strings.Contains,
        "replace":   strings.ReplaceAll,
    }

    return &Renderer{
        funcMap: funcMap,
    }
}

func (r *Renderer) Render(tmplPath string, ctx Context) (string, error) {
    tmpl, err := template.New(filepath.Base(tmplPath)).
        Funcs(r.funcMap).
        ParseFiles(tmplPath)

    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, ctx); err != nil {
        return "", err
    }

    return buf.String(), nil
}
```

### Configuration Management

```go
// internal/config/config.go
package config

type Config struct {
    Version          string
    TemplateRegistry []RegistrySource
    Defaults         Defaults
    ProjectDefaults  ProjectDefaults
    Behavior         Behavior
}

type Defaults struct {
    Language   string
    Docker     bool
    CIProvider string
    GitInit    bool

    Python LanguageDefaults
    NodeJS LanguageDefaults
    Kotlin LanguageDefaults
}

type LanguageDefaults struct {
    PackageManager string
    MinVersion     string
    Formatter      string
}

type ProjectDefaults struct {
    Author       string
    Email        string
    License      string
    Organization string
}

type Behavior struct {
    Interactive      bool
    ValidationLevel  string
    AutoGitInit      bool
    AutoInstallDeps  string  // yes, no, ask
}

// Load loads config from ~/.config/devinit/config.yaml
func Load() (*Config, error) {
    // Load and parse config
}

// Save saves config
func (c *Config) Save() error {
    // Save config
}
```

---

## Implementation Details

### MVP Scope

**Included in MVP**:
- Languages: Python (FastAPI), Node.js (Express)
- Features:
  - Docker support
  - Basic README
  - .env.example
  - .gitignore
  - Healthcheck endpoint
  - Basic validation
- CLI modes: Interactive and non-interactive
- Single binary distribution

**Post-MVP**:
- Additional languages (Kotlin, Java, Go, Rust)
- Additional frameworks per language
- Plugin system
- Template registry
- Project upgrades
- Advanced customization

### Development Phases

**Phase 1: Foundation** (Week 1-2)
```
- Project setup
- CLI framework (Cobra + Viper)
- Template loader
- Basic renderer
- Unit tests
```

**Phase 2: Core Features** (Week 3-4)
```
- Python/FastAPI template
- Node.js/Express template
- System validation
- Interactive prompts
- E2E tests
```

**Phase 3: Polish** (Week 5)
```
- Error handling
- User documentation
- CI/CD pipeline
- Release automation
```

**Phase 4: Distribution** (Week 6)
```
- Binary releases (Linux, macOS)
- Installation documentation
- Example projects
- Public release
```

### Technology Stack

**Core**:
- Language: Go 1.21+
- CLI Framework: Cobra
- Configuration: Viper
- Templates: Go text/template
- Interactive UI: Bubble Tea (optional)

**Testing**:
- Testing: Go testing package
- Assertions: testify
- Coverage: go-coverage
- E2E: Docker, HTTP testing

**CI/CD**:
- GitHub Actions
- GoReleaser (binary releases)
- Semantic versioning

**Distribution**:
- GitHub Releases (binaries)
- Homebrew (macOS)
- APT/YUM repos (Linux) - future
- Docker image - future

---

## Quality Assurance

### Code Quality

**Standards**:
- Go formatting: gofmt
- Linting: golangci-lint
- Code coverage: minimum 80%
- Documentation: godoc for all public APIs

**Pre-commit Hooks**:
```bash
#!/bin/bash
# .git/hooks/pre-commit

# Format
gofmt -w .

# Lint
golangci-lint run

# Test
go test ./...

# Validate templates
make validate-templates
```

### Template Quality

**Requirements**:
- All templates must have template.yaml
- All templates must pass validation
- All templates must have E2E test
- All generated projects must build successfully
- All generated projects must have working healthcheck

**Validation Checklist**:
```
□ template.yaml exists and is valid
□ All referenced files exist
□ All template variables are defined
□ No syntax errors in .tmpl files
□ Dockerfile builds successfully
□ Dependencies install successfully
□ Health endpoint returns 200
□ README has all required sections
□ Tests pass (if included)
```

### Performance Targets

```
CLI startup:          < 200ms
Project generation:   < 5s (small), < 30s (full with Docker)
Binary size:          < 10MB
Memory usage:         < 50MB during generation
Template validation:  < 1s per template
```

---

## Roadmap

### MVP (v1.0.0) - 6 weeks

**Deliverables**:
- CLI with core commands (new, validate, doctor, version)
- Python/FastAPI template
- Node.js/Express template
- Docker support
- Interactive and non-interactive modes
- System validation
- Binary releases (Linux, macOS)
- Documentation

**Success Metrics**:
- Generate working project in < 60 seconds
- Zero manual edits to run locally
- 100% template test coverage
- Positive feedback from 5+ beta users

### v1.1.0 - Post-MVP (Week 7-8)

**Features**:
- Add command (add components to existing projects)
- Shared components (CI templates)
- Template upgrade detection
- Improved error messages

### v1.2.0 - Enhancement (Week 9-12)

**Features**:
- Additional language: Kotlin/Spring Boot
- GitHub Actions CI template
- Prometheus monitoring template
- Project upgrade command
- Configuration file support

### v2.0.0 - Extensibility (Month 4+)

**Features**:
- Plugin system
- Template registry
- Custom template support
- Cloud provider integrations
- Advanced customization

---

## Appendix

### Example Template: Python/FastAPI

See: `templates/python/fastapi/` for complete example

### Example Generated Project

```
my-service/
├── .devinit.yaml
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── README.md
├── pyproject.toml
├── src/
│   ├── __init__.py
│   └── main.py
└── tests/
    ├── __init__.py
    └── test_main.py
```

### Contributing Guidelines

See: `CONTRIBUTING.md` (to be created)

### Release Process

See: `RELEASING.md` (to be created)

---

**Document Version**: 1.0.0
**Last Review**: 2025-12-25
**Next Review**: 2026-01-25
**Status**: Approved for implementation
