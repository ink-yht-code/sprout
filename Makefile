.PHONY: all build clean test lint fmt install help

all: build

build:
	@echo "Building sprout..."
	cd sprout && go build -v ./...
	@echo "Building sproutx..."
	cd sproutx && go build -v ./...
	@echo "Building sprout-gen..."
	cd sprout-gen && go build -v ./...
	@echo "Building sprout-registry..."
	cd sprout-registry && go build -v ./...
	@echo "Build complete!"

build-sprout:
	@echo "Building sprout..."
	cd sprout && go build -v ./...

build-sproutx:
	@echo "Building sproutx..."
	cd sproutx && go build -v ./...

build-sprout-gen:
	@echo "Building sprout-gen..."
	cd sprout-gen && go build -v ./...

build-sprout-registry:
	@echo "Building sprout-registry..."
	cd sprout-registry && go build -v ./...

clean:
	@echo "Cleaning..."
	go clean ./...
	cd sprout && go clean ./...
	cd sproutx && go clean ./...
	cd sprout-gen && go clean ./...
	cd sprout-registry && go clean ./...
	@echo "Clean complete!"

test:
	@echo "Running tests..."
	go test ./...
	cd sprout && go test ./...
	cd sproutx && go test ./...
	cd sprout-gen && go test ./...
	cd sprout-registry && go test ./...
	@echo "Tests complete!"

test-sprout:
	@echo "Running sprout tests..."
	cd sprout && go test ./...

test-sproutx:
	@echo "Running sproutx tests..."
	cd sproutx && go test ./...

test-sprout-gen:
	@echo "Running sprout-gen tests..."
	cd sprout-gen && go test ./...

test-sprout-registry:
	@echo "Running sprout-registry tests..."
	cd sprout-registry && go test ./...

lint:
	@echo "Running linter..."
	golangci-lint run ./...
	cd sprout && golangci-lint run ./...
	cd sproutx && golangci-lint run ./...
	cd sprout-gen && golangci-lint run ./...
	cd sprout-registry && golangci-lint run ./...
	@echo "Lint complete!"

fmt:
	@echo "Formatting code..."
	go fmt ./...
	cd sprout && go fmt ./...
	cd sproutx && go fmt ./...
	cd sprout-gen && go fmt ./...
	cd sprout-registry && go fmt ./...
	@echo "Format complete!"

install:
	@echo "Installing sprout-gen..."
	cd sprout-gen && go install ./...
	@echo "Install complete!"

run-registry:
	@echo "Running sprout-registry..."
	cd sprout-registry && go run cmd/main.go

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	cd sprout && go mod download && go mod tidy
	cd sproutx && go mod download && go mod tidy
	cd sprout-gen && go mod download && go mod tidy
	cd sprout-registry && go mod download && go mod tidy
	@echo "Dependencies updated!"

help:
	@echo "Sprout Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all                  Build all modules (default)"
	@echo "  build                Build all modules"
	@echo "  build-sprout        Build sprout module"
	@echo "  build-sproutx       Build sproutx module"
	@echo "  build-sprout-gen    Build sprout-gen module"
	@echo "  build-sprout-registry Build sprout-registry module"
	@echo "  clean                Clean all build artifacts"
	@echo "  test                 Run all tests"
	@echo "  test-sprout         Run sprout tests"
	@echo "  test-sproutx        Run sproutx tests"
	@echo "  test-sprout-gen     Run sprout-gen tests"
	@echo "  test-sprout-registry Run sprout-registry tests"
	@echo "  lint                 Run linter on all modules"
	@echo "  fmt                  Format all code"
	@echo "  install              Install sprout-gen"
	@echo "  run-registry         Run sprout-registry"
	@echo "  deps                 Download and tidy dependencies"
	@echo "  help                 Show this help message"
