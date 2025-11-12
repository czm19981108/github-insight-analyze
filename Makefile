.PHONY: all build run clean test help install deps

# Variables
BINARY_NAME=notifier
BUILD_DIR=bin
MAIN_PATH=./cmd/notifier
GO=go

# Detect OS for binary extension and color support
# 用于跨平台兼容 判断二进制文件要不要加后缀
# 注意 后续本项目的make命令使用都在git bash中执行
ifeq ($(OS),Windows_NT)
	BINARY_EXT=.exe
	# Windows: use printf with ANSI codes (works in modern Windows Terminal, Git Bash)
	ECHO=@printf
	COLOR_RESET=\033[0m
	COLOR_BOLD=\033[1m
	COLOR_GREEN=\033[32m
	COLOR_YELLOW=\033[33m
else
	BINARY_EXT=
	# Unix: use printf with ANSI codes
	ECHO=@printf
	COLOR_RESET=\033[0m
	COLOR_BOLD=\033[1m
	COLOR_GREEN=\033[32m
	COLOR_YELLOW=\033[33m
endif

all: deps build ## Build the application

help: ## Show this help message
	@echo '$(COLOR_BOLD)Available targets:$(COLOR_RESET)'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-15s$(COLOR_RESET) %s\n", $$1, $$2}'

deps: ## Download dependencies
	$(ECHO) "$(COLOR_YELLOW)Downloading dependencies...$(COLOR_RESET)\n"
	$(GO) mod download
	$(GO) mod verify

build: clean-build ## Build the application
	$(ECHO) "$(COLOR_YELLOW)Building $(BINARY_NAME)...$(COLOR_RESET)\n"
	@mkdir -p $(BUILD_DIR) 2>/dev/null || mkdir $(BUILD_DIR) 2>nul || true
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)$(BINARY_EXT) $(MAIN_PATH)
	$(ECHO) "$(COLOR_GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(BINARY_EXT)$(COLOR_RESET)\n"

clean-build: ## Clean build directory before building
	$(ECHO) "$(COLOR_YELLOW)Cleaning old build artifacts...$(COLOR_RESET)\n"
	@$(GO) clean -cache 2>/dev/null || true
	@rm -rf $(BUILD_DIR) 2>/dev/null || rmdir /s /q $(BUILD_DIR) 2>nul || true

install: ## Install the application
	$(ECHO) "$(COLOR_YELLOW)Installing $(BINARY_NAME)...$(COLOR_RESET)\n"
	$(GO) install $(MAIN_PATH)
	$(ECHO) "$(COLOR_GREEN)Installation complete$(COLOR_RESET)\n"

run: build ## Build and run the application
	$(ECHO) "$(COLOR_YELLOW)Running $(BINARY_NAME)...$(COLOR_RESET)\n"
	./$(BUILD_DIR)/$(BINARY_NAME)$(BINARY_EXT)

run-config: build ## Run with config file
	$(ECHO) "$(COLOR_YELLOW)Running $(BINARY_NAME) with config file...$(COLOR_RESET)\n"
	./$(BUILD_DIR)/$(BINARY_NAME)$(BINARY_EXT) -config configs/config.yaml

test: ## Run tests
	$(ECHO) "$(COLOR_YELLOW)Running tests...$(COLOR_RESET)\n"
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	$(ECHO) "$(COLOR_YELLOW)Running tests with coverage...$(COLOR_RESET)\n"
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	$(ECHO) "$(COLOR_GREEN)Coverage report generated: coverage.html$(COLOR_RESET)\n"

clean: ## Clean build artifacts
	$(ECHO) "$(COLOR_YELLOW)Cleaning...$(COLOR_RESET)\n"
	@rm -rf $(BUILD_DIR) 2>/dev/null || rmdir /s /q $(BUILD_DIR) 2>nul || true
	@rm -f coverage.out coverage.html 2>/dev/null || del /f coverage.out coverage.html 2>nul || true
	@$(GO) clean
	$(ECHO) "$(COLOR_GREEN)Clean complete$(COLOR_RESET)\n"

fmt: ## Format code
	$(ECHO) "$(COLOR_YELLOW)Formatting code...$(COLOR_RESET)\n"
	$(GO) fmt ./...
	$(ECHO) "$(COLOR_GREEN)Format complete$(COLOR_RESET)\n"

lint: ## Run linter
	$(ECHO) "$(COLOR_YELLOW)Running linter...$(COLOR_RESET)\n"
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

vet: ## Run go vet
	$(ECHO) "$(COLOR_YELLOW)Running go vet...$(COLOR_RESET)\n"
	$(GO) vet ./...

check: fmt vet test ## Run all checks (format, vet, test)

version: ## Show version
	@$(BUILD_DIR)/$(BINARY_NAME) -version 2>/dev/null || echo "Build the application first: make build"

docker-build: ## Build Docker image
	$(ECHO) "$(COLOR_YELLOW)Building Docker image...$(COLOR_RESET)\n"
	docker build -t ossinsight-notifier:latest .

docker-run: ## Run Docker container
	$(ECHO) "$(COLOR_YELLOW)Running Docker container...$(COLOR_RESET)\n"
	docker run --env-file .env ossinsight-notifier:latest

setup-env: ## Create .env file from example
# 缺少env时会基于example创建
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(COLOR_GREEN).env file created from .env.example$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Please edit .env with your configuration$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW).env file already exists$(COLOR_RESET)"; \
	fi

setup-config: ## Create config.yaml from example
# 缺少config时会基于example创建
	@if [ ! -f configs/config.yaml ]; then \
		cp configs/config.example.yaml configs/config.yaml; \
		echo "$(COLOR_GREEN)config.yaml created from config.example.yaml$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Please edit configs/config.yaml with your configuration$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)configs/config.yaml already exists$(COLOR_RESET)"; \
	fi
