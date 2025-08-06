.PHONY: init run migrate tidy install-air install-swag swagger-gen build

init: tidy install-air install-swag
	@echo "Project initialized successfully"

tidy:
	@echo "Tidying go modules..."
	go mod tidy

install-air:
	@echo "Checking if air is installed..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)

install-swag:
	@echo "Checking if swag is installed..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)

swagger-gen:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go -o docs/

build: tidy swagger-gen
	@echo "Building the application..."
	go build -o bin/server ./cmd/server

migrate:
	@echo "Running database migrations..."
	migrate -path migrations -database "postgres://pcrm_user:pcrm_password@localhost:5445/pcrm?sslmode=disable" up

run: init migrate swagger-gen
	@echo "Starting application with hot reload..."
	air