# devinit

Multi-language project scaffolding CLI tool that creates production-ready projects with standardized structure, Docker support, and best practices built-in.

## Features

- Production-ready templates for multiple languages and frameworks
- Docker and Docker Compose support out of the box
- Consistent project structure across different stacks
- Built-in healthcheck endpoints
- Test setup included
- Interactive and non-interactive modes

## Quick Start

### Installation

Build from source:

```bash
git clone https://github.com/renan-dev/devinit
cd devinit
make build
```

The binary will be available at `bin/devinit`.

### Usage

Create a new Python FastAPI project:

```bash
devinit new my-api --lang python --framework fastapi
```

With database support:

```bash
devinit new my-api \
  --lang python \
  --framework fastapi \
  --database postgres
```

### Available Templates

```bash
devinit templates list
```

Current templates:
- `python/fastapi` - FastAPI web framework with async support

### Commands

```bash
# Create new project
devinit new <name> --lang <language> --framework <framework>

# List available templates
devinit templates list

# Show template details
devinit templates show <template>

# Validate all templates
devinit templates validate

# Check system requirements
devinit doctor

# Validate existing project
devinit validate
```

## Development

### Prerequisites

- Go 1.21+
- Make

### Build

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Validate templates
make validate-templates

# Format code
make fmt
```

### Project Structure

```
devinit/
├── cmd/
│   └── devinit/          # CLI entrypoint
├── internal/
│   ├── generator/        # Project generator
│   ├── template/         # Template engine
│   ├── config/           # Configuration
│   ├── prompt/           # Interactive prompts
│   └── validator/        # Validation logic
├── templates/            # Project templates
│   ├── python/
│   │   └── fastapi/
│   └── shared/           # Shared components
├── test/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── docs/
└── Makefile
```

## Architecture

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed architecture documentation.

See [DECISIONS.md](./DECISIONS.md) for quick reference guide.

## Creating Templates

Templates are located in `templates/<language>/<framework>/`.

Each template must have:
- `template.yaml` - Template metadata and configuration
- `files/` - Directory containing template files

Example template structure:

```
templates/python/fastapi/
├── template.yaml
└── files/
    ├── main.py.tmpl
    ├── Dockerfile
    ├── README.md.tmpl
    └── ...
```

Files with `.tmpl` extension are processed as Go templates. Other files are copied as-is.

## Examples

### Create Python FastAPI project with PostgreSQL

```bash
devinit new user-service \
  --lang python \
  --framework fastapi \
  --database postgres \
  --tests

cd user-service
docker compose up
```

The API will be available at `http://localhost:8000` with:
- Swagger docs at `/docs`
- Health check at `/health`
- Database health check at `/db-health`

### Dry run to preview files

```bash
devinit new my-api \
  --lang python \
  --framework fastapi \
  --dry-run
```

## Roadmap

### v1.0.0 (MVP) - Current
- [x] Core CLI framework
- [x] Template engine
- [x] Python/FastAPI template
- [x] Docker support
- [ ] System validation
- [ ] E2E tests

### v1.1.0
- [ ] Node.js/Express template
- [ ] Add command (add components to existing projects)
- [ ] CI templates (GitHub Actions)

### v1.2.0
- [ ] Kotlin/Spring Boot template
- [ ] Project upgrade command
- [ ] Configuration file support

### v2.0.0
- [ ] Plugin system
- [ ] Template registry
- [ ] Custom template support

## License

MIT

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## Support

- Issues: https://github.com/renan-dev/devinit/issues
- Documentation: ./docs/
