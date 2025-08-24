# ConfigCraft Build Automation

.DEFAULT_GOAL := help

# Variables
BINARY_NAME=configcraft
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-s -w -H windowsgui -X configcraft/internal/version.Version=$(VERSION)"

## help: Show this help message
.PHONY: help
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## deps: Install and verify dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	go mod verify
	@echo "Dependencies installed successfully!"

## build: Build the application
.PHONY: build
build: deps
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe main.go
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME).exe"

## dev: Run in development mode
.PHONY: dev
dev:
	@echo "Running in development mode..."
	go run main.go

## cli: Run CLI version
.PHONY: cli
cli:
	@echo "Running CLI version..."
	cd cmd && go run cli.go

## test: Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

## clean: Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -rf $(BUILD_DIR)/$(BINARY_NAME).exe
	@echo "Clean completed!"

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted!"

## lint: Run linter (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run
	@echo "Linting completed!"

## mod: Update dependencies
.PHONY: mod
mod:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies updated!"

## install: Install the application globally
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) globally..."
	cp $(BUILD_DIR)/$(BINARY_NAME).exe $(GOPATH)/bin/
	@echo "Installation completed!"

## release: Create release build with version info
.PHONY: release
release: clean
	@echo "Creating release build $(VERSION)..."
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION).exe main.go
	@echo "Release build completed: $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION).exe"