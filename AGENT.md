# Personal CRM Project

## Architecture Overview

This is a personal CRM application backend app built with Golang. App is dedicated to in-person interaction management.

### Technology Stack
- **Backend**: Go 1.23.x with built-in `net/http` routing (Go 1.22+ features)
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
│   ├── dto/              # API contracts (request/response structures)
│   ├── handlers/         # HTTP handlers (orchestration layer)
│   ├── mappers/          # Data transformation between layers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Domain models
│   ├── repository/       # Database operations
│   ├── services/         # Business logic
│   └── validators/       # Input validation logic
├── migrations/           # Database migrations
├── docker-compose.yml   # Development infrastructure
└── config.yml          # Application configuration
```

## Architecture Principles

### Separation of Concerns
- **DTOs**: API contracts only, no business logic
- **Validators**: Pure validation functions, no data transformation
- **Mappers**: Data transformation between layers only
- **Services**: Business logic only, no DTOs or HTTP concerns
- **Handlers**: HTTP orchestration only, delegate to other layers

### Layer Boundaries
- No circular dependencies between layers
- Each package has single responsibility

### Domain-Driven Design
- Models represent business domain
- Services contain business rules and logic
- DTOs are transport contracts, not domain objects
- Clear separation between API and domain concerns

### Minimal Dependencies
- Prefer standard library where possible

## Code Style
- Self-documented code: Write code that explains itself without extra comments. If comment is required for better explanation code must be extracted into separate function and then function comment is added. Inline comments are NOT ALLOWED untill explicitly requested
- Method and function bodies may not contain comments
- Method bodies may not contain blank lines
- Expressive naming: Make the naming in code self-documenting. The name should reflect the full functionality of the function / method
- Explicit contracts: Function ALWAYS MUST do only what it said to do, no implicit actions allowed, no side-effects
- Error and log messages should not end with a period
- Prefer explicit error handling over panics
- All database changes must go through migrations

## Development Guidelines

**IMPORTANT**: When developing, assume the application is already running with air hot reload. NEVER run the app manually - only use `make build` or `make test`.

### Allowed Make Commands

- **`make init`**: Initialize project (tidy modules, install air, install swag)
- **`make build`**: Build production binary (tidy + swagger-gen + build to bin/server)
- **`make test`**: Run all tests
- **`make migrate`**: Run database migrations
- **`make swagger-gen`**: Generate Swagger documentation from code annotations
- **`make tidy`**: Clean up Go modules

### Development Workflow

- Application runs with air hot reload during development
- Use `make build` to verify compilation without running
- Use `make test` to run tests
- ALWAYS check Makefile for available commands before creating custom scripts. You must create custom scripts ONLY with explicit user concent

## Configuration

The application uses YAML configuration with environment-specific overrides located in /config.yaml

## Git Commit Guidelines
- Commits: short imperative subject (e.g., "add swagger autogen", "fix bug"), keep context clear.
- Link related issues; update `docs/` (Swagger) and `README.md` if endpoints or setup change.
- When adding/modifying endpoints, run `make swagger-gen` to update API documentation.
- Ensure `make build` and `make test` succeed, and migrations are idempotent.
