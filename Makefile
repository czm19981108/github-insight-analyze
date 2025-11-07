.PHONY: all build run clean test help install deps

# Variables
BINARY_NAME=notifier
BUILD_DIR=bin
MAIN_PATH=./cmd/notifier
GO=go

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m

all: deps build ## Build the application

help: ## Show this help message
	@echo '$(COLOR_BOLD)Available targets:$(COLOR_RESET)'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-15s$(COLOR_RESET) %s\n", $$1, $$2}'

deps: ## Download dependencies
	@echo "$(COLOR_YELLOW)Downloading dependencies...$(COLOR_RESET)"
	$(GO) mod download
	$(GO) mod verify

build: ## Build the application
	@echo "$(COLOR_YELLOW)Building $(BINARY_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(COLOR_GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

install: ## Install the application
	@echo "$(COLOR_YELLOW)Installing $(BINARY_NAME)...$(COLOR_RESET)"
	$(GO) install $(MAIN_PATH)
	@echo "$(COLOR_GREEN)Installation complete$(COLOR_RESET)"

run: build ## Build and run the application
	@echo "$(COLOR_YELLOW)Running $(BINARY_NAME)...$(COLOR_RESET)"
	./$(BUILD_DIR)/$(BINARY_NAME)

run-config: build ## Run with config file
	@echo "$(COLOR_YELLOW)Running $(BINARY_NAME) with config file...$(COLOR_RESET)"
	./$(BUILD_DIR)/$(BINARY_NAME) -config configs/config.yaml

test: ## Run tests
	@echo "$(COLOR_YELLOW)Running tests...$(COLOR_RESET)"
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(COLOR_YELLOW)Running tests with coverage...$(COLOR_RESET)"
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)Coverage report generated: coverage.html$(COLOR_RESET)"

clean: ## Clean build artifacts
	@echo "$(COLOR_YELLOW)Cleaning...$(COLOR_RESET)"
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	$(GO) clean
	@echo "$(COLOR_GREEN)Clean complete$(COLOR_RESET)"

fmt: ## Format code
	@echo "$(COLOR_YELLOW)Formatting code...$(COLOR_RESET)"
	$(GO) fmt ./...
	@echo "$(COLOR_GREEN)Format complete$(COLOR_RESET)"

lint: ## Run linter
	@echo "$(COLOR_YELLOW)Running linter...$(COLOR_RESET)"
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

vet: ## Run go vet
	@echo "$(COLOR_YELLOW)Running go vet...$(COLOR_RESET)"
	$(GO) vet ./...

check: fmt vet test ## Run all checks (format, vet, test)

version: ## Show version
	@$(BUILD_DIR)/$(BINARY_NAME) -version 2>/dev/null || echo "Build the application first: make build"

docker-build: ## Build Docker image
	@echo "$(COLOR_YELLOW)Building Docker image...$(COLOR_RESET)"
	docker build -t ossinsight-notifier:latest .

docker-run: ## Run Docker container
	@echo "$(COLOR_YELLOW)Running Docker container...$(COLOR_RESET)"
	docker run --env-file .env ossinsight-notifier:latest

setup-env: ## Create .env file from example
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(COLOR_GREEN).env file created from .env.example$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Please edit .env with your configuration$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW).env file already exists$(COLOR_RESET)"; \
	fi

setup-config: ## Create config.yaml from example
	@if [ ! -f configs/config.yaml ]; then \
		cp configs/config.example.yaml configs/config.yaml; \
		echo "$(COLOR_GREEN)config.yaml created from config.example.yaml$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Please edit configs/config.yaml with your configuration$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)configs/config.yaml already exists$(COLOR_RESET)"; \
	fi
