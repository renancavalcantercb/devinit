# Quick Decision Reference

**Last Updated**: 2025-12-25

## Critical Decisions

### Tech Stack
```
Language:     Go 1.21+
CLI:          Cobra + Viper
Templates:    Go text/template
Testing:      Go testing + testify
Distribution: GoReleaser
```

### Repository Structure
```
Monorepo: CLI + Templates together
Reason: Atomic releases, simpler MVP
Migration: Split when templates > 50MB
```

### Versioning
```
Semver: MAJOR.MINOR.PATCH
Major version must match between CLI and templates
Breaking changes only in MAJOR
```

### Validation
```
Three tiers:
1. Pre-flight:  System requirements check
2. Runtime:     Template variable validation
3. Post-gen:    Generated project validation

Default: Basic (warns but continues)
Flags:   --strict, --no-validate
```

---

## Implementation Priorities

### MVP Must-Have
```
✓ Languages:      Python, Node.js
✓ Frameworks:     FastAPI, Express
✓ Features:       Docker, README, .env.example, .gitignore
✓ CLI modes:      Interactive + non-interactive
✓ Validation:     System requirements
✓ Distribution:   Binary (Linux, macOS)
```

### Post-MVP
```
○ Languages:      Kotlin, Java, Go, Rust
○ Frameworks:     Flask, NestJS, Spring Boot
○ Features:       CI templates, monitoring, upgrades
○ Advanced:       Plugins, template registry
```

---

## Template Standards

### Required Files
```
template.yaml       # Metadata
files/              # Template files
  ├── README.md.tmpl
  ├── .gitignore
  ├── .env.example
  └── (language-specific)
```

### template.yaml Required Sections
```yaml
version: "1.0.0"
name: "Template Name"
language: python
framework: fastapi
requirements:
  system: [...]
variables: {...}
files: [...]
```

### Healthcheck Requirement
```
All templates must include a healthcheck endpoint
HTTP endpoint returning 200 on /health or /healthz
```

---

## File Naming Conventions

### Template Files
```
.tmpl extension:  Files with Go template syntax
No extension:     Static files (copied as-is)

Examples:
  main.py.tmpl      → Rendered with variables
  Dockerfile        → Copied as-is
  README.md.tmpl    → Rendered with variables
```

### Generated Projects
```
Pattern: ^[a-z][a-z0-9-]*$
Valid:   my-service, api-gateway, user-service
Invalid: MyService, my_service, 123-service
```

---

## CLI Behavior

### Default Mode
```bash
devinit new              # Interactive wizard
devinit new api foo      # Interactive with project name
```

### Non-Interactive
```bash
devinit new api foo --lang python --framework fastapi --docker
```

### Validation Modes
```bash
--no-validate    # Skip all validation
(default)        # Basic validation (warn but continue)
--strict         # Strict validation (error and stop)
```

### Dry Run
```bash
devinit new api foo --dry-run    # Show what would be done
```

---

## Error Handling

### Error Levels
```
error:   Block execution, return non-zero
warn:    Show warning, continue
ignore:  Silently continue
```

### User-Facing Errors
```
Format:
  Error: <what happened>
  Cause: <why it happened>
  Fix:   <how to resolve>

Example:
  Error: Docker not found
  Cause: Docker is required when --docker flag is used
  Fix:   Install Docker: https://docs.docker.com/install
```

---

## Testing Requirements

### Per Template
```
1. Unit test:       Template metadata validation
2. Render test:     All files render without error
3. E2E test:        Generated project builds and runs
4. Health test:     Healthcheck endpoint works
```

### E2E Test Flow
```
1. Generate project
2. Validate file structure
3. Build (Docker or native)
4. Run service
5. Test healthcheck
6. Cleanup
```

### Coverage Target
```
Overall:    80%
Core pkg:   90%
Templates:  100% (must have E2E test)
```

---

## Configuration

### Global Config Location
```
Linux:   ~/.config/devinit/config.yaml
macOS:   ~/.config/devinit/config.yaml
Windows: %APPDATA%\devinit\config.yaml
```

### Project Metadata Location
```
.devinit.yaml in project root
```

### Environment Variables
```
DEVINIT_CONFIG:     Override config location
DEVINIT_NO_COLOR:   Disable colored output
DEVINIT_LOG_LEVEL:  Set log level (debug, info, warn, error)
```

---

## Go Template Functions

### Built-in Available
```
lower:      strings.ToLower
upper:      strings.ToUpper
contains:   strings.Contains
replace:    strings.ReplaceAll
```

### Custom Functions (to implement)
```
snake:      toSnakeCase (my_project)
camel:      toCamelCase (myProject)
pascal:     toPascalCase (MyProject)
kebab:      toKebabCase (my-project)
```

### Usage Example
```
{{ .ProjectName | lower }}
{{ .ProjectName | snake }}
{{ .Description | replace "." "" }}
```

---

## Directory Structure (Generated Projects)

### Python/FastAPI
```
project-name/
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

### Node.js/Express
```
project-name/
├── .devinit.yaml
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── README.md
├── package.json
├── src/
│   └── index.js
└── tests/
    └── index.test.js
```

---

## Build & Release

### Build Commands
```bash
make build              # Build for current platform
make build-all          # Build for all platforms
make test               # Run all tests
make test-unit          # Unit tests only
make test-e2e           # E2E tests only
make validate-templates # Validate all templates
```

### Release Process
```
1. Update version in code
2. Update CHANGELOG.md
3. Tag: git tag v1.0.0
4. Push: git push --tags
5. CI builds and releases binaries
```

### Binary Naming
```
devinit-{version}-{os}-{arch}

Examples:
  devinit-1.0.0-linux-amd64
  devinit-1.0.0-darwin-arm64
  devinit-1.0.0-windows-amd64.exe
```

---

## Dependency Management

### Go Dependencies
```
Minimal dependencies:
  - cobra (CLI)
  - viper (config)
  - yaml.v3 (YAML parsing)
  - testify (testing only)

Avoid heavy dependencies
```

### Template Dependencies
```
Declared in template.yaml requirements section
Validated at generation time
User installs after project creation
```

---

## Documentation Requirements

### README Sections
```
Required:
  - Project description
  - Prerequisites
  - Installation
  - Running locally
  - Running with Docker
  - Health check
  - Project structure

Optional:
  - Deployment
  - Testing
  - Contributing
  - License
```

### Code Documentation
```
All public functions:   godoc comments
All interfaces:         godoc comments
Complex logic:          inline comments
Examples:               godoc examples
```

---

## Security Considerations

### Template Variables
```
Always validate:
  - Project names (alphanumeric + hyphens)
  - File paths (no path traversal)
  - Command injection (no shell execution)

Never trust user input without validation
```

### Generated Code
```
Do not include:
  - Hardcoded credentials
  - API keys
  - Secrets

Use .env.example with placeholder values
```

### Dependencies
```
Review all dependencies for:
  - Known vulnerabilities
  - License compatibility
  - Maintenance status
```

---

## Performance Guidelines

### Template Rendering
```
Target: < 1s for average template
Optimize: Cache parsed templates
Avoid:    Excessive file I/O
```

### Project Generation
```
Target: < 5s for small project, < 30s with Docker
Optimize: Parallel file operations where possible
Avoid:    Sequential operations that could be parallel
```

### Binary Size
```
Target:     < 10MB
Technique:  Strip debug symbols in release
Embed:      Templates using go:embed
```

---

## Common Pitfalls to Avoid

### Template Design
```
✗ Don't: Use language-specific logic in templates
✓ Do:    Keep templates simple, logic in generator

✗ Don't: Hardcode values
✓ Do:    Use variables for everything configurable

✗ Don't: Copy-paste between templates
✓ Do:    Use shared components
```

### Code Organization
```
✗ Don't: Put logic in cmd/
✓ Do:    Keep cmd/ thin, logic in internal/

✗ Don't: Export everything
✓ Do:    Keep APIs minimal, unexported by default

✗ Don't: Global state
✓ Do:    Pass context explicitly
```

### Error Handling
```
✗ Don't: Panic in library code
✓ Do:    Return errors

✗ Don't: Ignore errors
✓ Do:    Handle or wrap with context

✗ Don't: Generic error messages
✓ Do:    Specific, actionable error messages
```

---

## Migration Path (Future)

### When to Split Repo
```
Triggers:
  - Templates > 50MB
  - Multiple teams maintaining different templates
  - Community contributing templates
  - Need independent template versioning

Migration:
  1. Create devinit-templates repo
  2. Move templates/ directory
  3. Update CI to fetch templates
  4. Version compatibility matrix
  5. Deprecation period for old CLI
```

### Plugin System (v2.0+)
```
After:
  - Core is stable
  - Clear plugin API surface
  - Use cases validated

Implementation:
  - Go plugin system OR
  - External binary plugins (safer)
```

---

## Support Matrix

### Operating Systems (MVP)
```
✓ Linux:   amd64, arm64
✓ macOS:   amd64 (Intel), arm64 (Apple Silicon)
○ Windows: Future (post-MVP)
```

### Go Version
```
Minimum: Go 1.21
Target:  Latest stable (1.21+)
```

### Docker Version
```
Minimum: 24.0 (if using Docker features)
Recommended: Latest stable
```

---

## Quick Commands Reference

### Development
```bash
# Setup
make setup

# Build
make build

# Test
make test
make test-unit
make test-integration
make test-e2e

# Validate
make validate-templates
make lint

# Run
./bin/devinit new

# Clean
make clean
```

### Usage
```bash
# Interactive
devinit new

# Quick start
devinit new api my-service --lang python

# Full options
devinit new api my-service \
  --lang python \
  --framework fastapi \
  --docker \
  --database postgres \
  --ci github

# Validate
devinit validate
devinit doctor

# Templates
devinit templates list
devinit templates show python/fastapi
```

---

## When in Doubt

### Ask These Questions
```
1. Is this needed for MVP?
   → If no, defer to post-MVP

2. Does this add complexity?
   → If yes, justify the benefit

3. Can this be done simpler?
   → Always prefer simpler

4. Is this tested?
   → If no, add tests

5. Is this documented?
   → If no, add docs
```

### Principles
```
1. Simple is better than complex
2. Explicit is better than implicit
3. Production-ready over feature-rich
4. User experience over developer convenience
5. Ship MVP fast, iterate based on feedback
```

---

**Document Version**: 1.0.0
**Companion to**: ARCHITECTURE.md
**Status**: Reference guide
