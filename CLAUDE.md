# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TheOnePager is a Go web application that serves as a homepage dashboard for virtual home management. It's built with:
- **Backend**: Go with Echo web framework
- **Frontend**: TailwindCSS with DaisyUI components
- **Architecture**: Standard Go project layout with internal/pkg separation
- **Hot Reload**: Air for development workflow
- **Testing**: Ginkgo BDD testing framework

## Common Development Commands

### Development Server
```bash
make run          # Start development server with hot reload (uses Air)
air               # Direct Air command for hot reload
```

### CSS Development
```bash
make css          # Watch and build TailwindCSS
npm run css       # Alternative CSS watch command
npm run css::prod # Production CSS build (minified)
```

### Testing and Quality
```bash
make test         # Run all tests with Ginkgo
ginkgo -r         # Direct Ginkgo command
make testr        # Auto-rerun tests on file changes
make lint         # Run golangci-lint with all rules enabled
```

### Database Operations
```bash
make sqlgen       # Generate SQL with sqlc
make sqlgenr      # Auto-regenerate SQL on file changes
make migrate      # Run database migrations
```

### Dependency Management
```bash
make update       # Update all Go dependencies
go mod tidy       # Clean up module dependencies
```

## Project Architecture

### Core Structure
- **main.go**: Entry point with CLI flag parsing (`--serve`, `--port`, `--config`)
- **src/run.go**: Application runner with graceful shutdown handling
- **src/internal/app/**: Core application logic
  - `app.go`: Main App struct and server setup
  - `routes.go`: Route definitions
  - `config.go`: Configuration parsing
  - `handle_*.go`: HTTP handlers
- **src/pkg/**: Shared packages (logger, renderer)
- **src/web/**: Static assets and templates

### Configuration
- **config.yaml**: Application configuration with dashboard items
- **Environment**: Uses envconfig for environment-based configuration
- **Templates**: HTML templates in `src/web/templates/`
- **Static Assets**: CSS/JS in `src/web/static/`

### Key Dependencies
- Echo v4 for HTTP server
- TailwindCSS + DaisyUI for styling
- Zerolog for structured logging
- Unrolled packages for security middleware
- Live reload support in development

### Development Workflow
1. Use `make run` for development with hot reload
2. CSS changes require `make css` in parallel
3. Tests run with `make test` or `make testr` for auto-rerun
4. Lint before commits with `make lint`

The application serves a configurable dashboard showing applications with icons and links, defined in config.yaml.