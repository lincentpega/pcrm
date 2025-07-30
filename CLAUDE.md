# Personal CRM Project

## Architecture Overview

This is a personal CRM application built with a simple, modern Go stack focused on server-side rendering and minimal dependencies.

### Technology Stack
- **Backend**: Go 1.23.x with built-in `net/http` routing (Go 1.22+ features)
- **Frontend**: HTMX for dynamic interactions
- **Templates**: `html/template` for server-side rendering
- **Database**: PostgreSQL with `sqlx` for type-safe queries
- **Migrations**: `golang-migrate` for database schema management
- **Middleware**: `alice` for middleware chaining
- **Configuration**: YAML-based configuration
- **Development**: Docker + Docker Compose

## Project Structure

```
pcrm/
├── cmd/
│   └── server/           # Application entrypoint
├── internal/
│   ├── config/           # Configuration management
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models
│   ├── repository/       # Database operations
│   └── templates/        # HTML templates
├── migrations/           # Database migrations
├── static/              # Static assets (CSS, JS, images)
├── docker-compose.yml   # Development infrastructure
├── Dockerfile           # Production container
└── config.yml          # Application configuration
```

## Getting Started

### Prerequisites
- Go 1.24.x or later
- Docker and Docker Compose
- golang-migrate CLI tool
- Air (for hot reload during development)

## Debugging Guidelines

When developing application instance of an app is already running with air, so there is no need to start an appp

When running some application-related commands first check Makefile content, there might be some command already pre-defined for you

## Key Design Decisions

- **No external router**: Using Go 1.22+ built-in routing improvements
- **Server-side rendering**: HTMX + html/template for dynamic UI without heavy JavaScript
- **Minimal dependencies**: Prefer standard library where possible
- **Simple structure**: Avoiding over-engineering patterns like DDD for this scale
- **Type safety**: Using sqlx for compile-time query validation

## Configuration

The application uses YAML configuration with environment-specific overrides located in /config.yaml

## Development Notes

- Use Docker Compose for consistent development environment
- All database changes must go through migrations
- Templates should be organized by feature/page
- HTMX responses should return HTML fragments when possible
- Use Alice for clean middleware chaining
- Prefer explicit error handling over panics

## Code Style Guidelines

- DO NOT ADD comments unless explicitly requested
- Keep code clean and self-documenting through clear naming
- Function and variable names should explain their purpose

## Git Commit Guidelines

- DO NOT add any mentions of Claude Code, AI assistance, or similar references in commit messages
- Keep commit messages focused on the technical changes and their purpose