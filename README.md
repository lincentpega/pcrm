# Personal CRM Backend API

A minimal personal CRM backend built with Go. It focuses on clean boundaries (DTOs, validators, mappers, repositories) and minimal dependencies while exposing a typed REST API with Swagger docs.

## Features

- **People**: Create, read, update, delete people
- **Contacts**: Store multiple contact methods and list available contact types
- **Connection Source**: Track how you met a person (meeting story, introducer)
- **Birth Date Info**: Store exact/partial birth date or approximate age
- **Conversations**: Log interactions with type, initiator, and notes; list conversation types
- **Pagination**: Basic pagination for people list
- **OpenAPI**: Swagger UI available under `/swagger`

## Technology Stack

- **Backend**: Go 1.23.x using `net/http` (1.22+ patterns)
- **Database**: PostgreSQL with `sqlx`
- **Migrations**: `golang-migrate`
- **Middleware**: `alice`
- **Configuration**: YAML (`config.yml`)
- **Development**: Docker + Docker Compose

## Quick Start

### Prerequisites

- Go 1.23.x
- Docker and Docker Compose
- `golang-migrate` CLI

### Setup

1) Start PostgreSQL:
```bash
docker compose up -d postgres
```

2) Apply migrations:
```bash
make migrate
```

3) Work on the codebase:
- App runs with air hot reload during development (assume it’s already running)
- Use `make build` to verify compilation and generate Swagger
- Use `make test` to run tests

### Make Commands

```bash
make init         # Install dev tools and tidy modules
make tidy         # Tidy go.mod/go.sum
make swagger-gen  # Generate Swagger docs (docs/)
make build        # Build production binary to bin/server
make test         # Run tests
make migrate      # Run DB migrations
```

## Project Structure

```
pcrm/
├── cmd/
│   └── server/            # Application entrypoint
├── internal/
│   ├── config/            # Configuration management
│   ├── dto/               # API contracts (request/response structures)
│   ├── handlers/          # HTTP handlers (orchestration layer)
│   ├── mappers/           # Data transformation between layers
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Domain models
│   ├── repository/        # Database operations
│   └── validators/        # Input validation logic
├── migrations/            # Database migrations
├── docker-compose.yml     # Development infrastructure
└── config.yml             # Application configuration
```

## API

- Swagger UI: `GET /swagger`
- People: `GET/POST /api/people`, `GET/PUT/DELETE /api/people/{id}`
- Contacts: `GET /api/people/{personId}/contacts`, `POST /api/people/{personId}/contacts`, `GET/PUT/DELETE /api/contacts/{id}`, `GET /api/contact-types`
- Connection Source: `GET/PUT/DELETE /api/people/{personId}/connection-source`
- Birth Date Info: `GET/PUT/DELETE /api/people/{personId}/birth-date-info`
- Conversations: `GET /api/people/{personId}/conversations`, `POST /api/people/{personId}/conversations`, `GET/PUT/DELETE /api/conversations/{id}`, `GET /api/conversation-types`

## Notes

- All database changes must go through migrations
- Error messages do not end with a period
- Prefer explicit error handling over panics

## License

MIT License - see LICENSE file for details
