.PHONY: run build test test-coverage clean wire mock migrate swagger swagger-clean swagger-init deps help

.DEFAULT_GOAL := help

APP_NAME = go-rest-api
BINARY_PATH = bin/$(APP_NAME)

GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOMOD = $(GOCMD) mod
GOCLEAN = $(GOCMD) clean

SWAG = swag
WIRE = wire
MOCKERY = mockery

MYSQL_USER ?= root
MYSQL_PASSWORD ?= "I1mme/&W0S8P"
MYSQL_HOST ?= 127.0.0.1
MYSQL_PORT ?= 3306
MYSQL_DBNAME ?= go_api_db
MYSQL_CMD = mysql -u $(MYSQL_USER) -p$(MYSQL_PASSWORD) -h $(MYSQL_HOST) -P $(MYSQL_PORT)

help:
	@echo "Usage: make [target]"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: Run the application

run-task-1:
	$(GOCMD) run cmd/person/main.go

run-task-2:
	$(GOCMD) run cmd/server/main.go
	
## build: Build the binary
build:
	$(GOBUILD) -o $(BINARY_PATH) cmd/server/main.go

## test: Run all unit tests
test:
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

## clean: Remove build artifacts and coverage files
clean:
	rm -rf bin/
	rm -f coverage.out
	$(GOCLEAN)

## wire: Generate dependency injection code using Google Wire
wire:
	cd wire && $(WIRE)

## mock: Generate mocks using Mockery
mock:
	$(MOCKERY) --all --dir=internal --output=mocks

## swagger-init: Generate Swagger documentation
swagger-init:
	$(SWAG) init -g cmd/server/main.go -o docs

## swagger-clean: Remove generated Swagger docs
swagger-clean:
	rm -rf docs/

## swagger: Regenerate Swagger docs (clean + init)
swagger:
	swag init -g cmd/server/main.go -o docs


## migrate: Create database and run migrations
migrate:
	@echo "Creating database if not exists..."
	$(MYSQL_CMD) -e "CREATE DATABASE IF NOT EXISTS $(MYSQL_DBNAME);"
	@echo "Running migrations..."
	$(MYSQL_CMD) $(MYSQL_DBNAME) < migrations/001_create_users_table.sql
	@echo "Migration completed."

## deps: Download and install project dependencies and tools
deps:
	$(GOMOD) download
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/vektra/mockery/v2@latest
	go install github.com/swaggo/swag/cmd/swag@latest