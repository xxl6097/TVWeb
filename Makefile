.PHONY: help build run test clean install

# Build variables
BINARY_NAME=TVWeb
BUILD_DIR=./build
DIST_DIR=./dist
MAIN_PATH=./cmd/app/main.go

# Git variables
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S')

# Linker flags
LDFLAGS=-ldflags "-X github.com/xxl6097/TVWeb/pkg/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/xxl6097/TVWeb/pkg/version.BuildTime=$(BUILD_TIME)"

help:
	@echo "TVWeb - Makefile targets:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build files"
	@echo "  make deps           - Download dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Lint code"

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

build: deps
	@echo "Building $(BINARY_NAME)..."
	mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

build-dist: deps
	@echo "Building for distribution..."
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Distribution build completed"

run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	go test -v -cover ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

lint:
	@echo "Linting code..."
	golint ./...
	go vet ./...

clean:
	@echo "Cleaning build files..."
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	go clean
	@echo "Clean completed"

install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "Installation completed"
