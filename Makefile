# Go parameters
APP_NAME = db-backup
BUILD_DIR = bin
GO_FILES = $(shell find . -type f -name '*.go')
VERSION = 1.0.0

# Docker parameters
DOCKER_IMAGE = db-backup:latest
DOCKER_CONTAINER = db-backup-test
REGISTRY = myregistry.com/myuser

.PHONY: all build run fmt lint test docker-build docker-run docker-push clean

## Build the Go binary
build: fmt $(GO_FILES)
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go
	@echo "✅ Build complete! Binary at $(BUILD_DIR)/$(APP_NAME)"

## Run the application
run:
	@echo "🚀 Running $(APP_NAME)..."
	@go run ./cmd/main.go

## Format Go source code
fmt:
	@echo "📝 Formatting code..."
	@go fmt ./...

## Lint Go code using golangci-lint
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run

## Run unit tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

## Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

## Run Docker container
docker-run:
	@echo "🚀 Running $(DOCKER_IMAGE) in a container..."
	@docker run --rm -p 3306:3306 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

## Push Docker image to registry
docker-push:
	@echo "📦 Pushing $(DOCKER_IMAGE) to $(REGISTRY)..."
	@docker tag $(DOCKER_IMAGE) $(REGISTRY)/$(DOCKER_IMAGE)
	@docker push $(REGISTRY)/$(DOCKER_IMAGE)

## Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete!"
