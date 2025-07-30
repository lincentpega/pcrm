# Personal CRM

A simple, modern personal CRM application built with Go and HTMX, focused on server-side rendering and minimal dependencies.

## Features

- **Person Management**: Add, view, edit, and delete personal contacts
- **Contact Information**: Store multiple contact methods (email, phone, social media)
- **Pagination**: Navigate through large contact lists efficiently
- **Responsive Design**: Works seamlessly on desktop and mobile devices
- **Real-time Updates**: Dynamic interactions powered by HTMX

## Technology Stack

- **Backend**: Go 1.23+ with built-in `net/http` routing
- **Frontend**: HTMX for dynamic interactions
- **Templates**: `html/template` for server-side rendering
- **Database**: PostgreSQL with `sqlx` for type-safe queries
- **Migrations**: `golang-migrate` for database schema management
- **Styling**: Custom CSS with responsive design
- **Development**: Docker + Docker Compose with hot reload

## Quick Start

### Prerequisites

- Go 1.24.x or later
- Docker and Docker Compose
- golang-migrate CLI tool

### Running the Application

1. Clone the repository:
   ```bash
   git clone https://github.com/lincentpega/pcrm.git
   cd pcrm
   ```

2. Start the development environment:
   ```bash
   make dev
   ```

3. The application will be available at `http://localhost:8080`

### Available Commands

```bash
make dev          # Start development environment with hot reload
make build        # Build the application
make test         # Run tests
make migrate-up   # Run database migrations
make migrate-down # Rollback database migrations
```

## Project Structure

```
pcrm/
├── cmd/server/           # Application entrypoint
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
└── config.yml          # Application configuration
```

## Key Design Decisions

- **Server-side rendering**: Uses HTMX + html/template for dynamic UI without heavy JavaScript
- **Minimal dependencies**: Prefers Go standard library where possible
- **Simple architecture**: Avoids over-engineering patterns for this scale
- **Type safety**: Uses sqlx for compile-time query validation
- **Clean styling**: Light black theme with responsive design

## Development

The application uses Docker Compose for a consistent development environment. Database changes must go through migrations, and templates are organized by feature/page.

## License

MIT License - see LICENSE file for details