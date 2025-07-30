.PHONY: init run migrate tidy install-air

init: tidy install-air
	@echo "Project initialized successfully"

tidy:
	@echo "Tidying go modules..."
	go mod tidy

install-air:
	@echo "Checking if air is installed..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)

migrate:
	@echo "Running database migrations..."
	migrate -path migrations -database "postgres://pcrm_user:pcrm_password@localhost:5445/pcrm?sslmode=disable" up

run: init migrate
	@echo "Starting application with hot reload..."
	air