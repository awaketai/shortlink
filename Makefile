# Makefile for shortlink service

# Variables
APP_NAME := shortlink
PID_FILE := $(APP_NAME).pid
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Default target
.PHONY: all
all: build

# Install dependencies
.PHONY: install
install:
	go mod tidy

# Build the application
.PHONY: build
build:
	go build -o $(APP_NAME) main.go

# Run the application
.PHONY: run
run:
	go run main.go

# Start the application in background
.PHONY: start
start: $(APP_NAME)
	@if [ -f $(PID_FILE) ]; then \
		echo "Service is already running with PID $$(cat $(PID_FILE))"; \
		exit 1; \
	fi
	@echo "Starting $(APP_NAME) service..."
	@nohup ./$(APP_NAME) > $(APP_NAME).log 2>&1 & echo $$! > $(PID_FILE)
	@echo "Service started with PID $$(cat $(PID_FILE))"

# Build target as dependency
$(APP_NAME):
	@go build -o $(APP_NAME) main.go

# Stop the application
.PHONY: stop
stop:
	@if [ ! -f $(PID_FILE) ]; then \
		echo "Service is not running"; \
		exit 1; \
	fi
	@echo "Stopping $(APP_NAME) service with PID $$(cat $(PID_FILE))..."
	@kill -TERM $$(cat $(PID_FILE))
	@sleep 2
	@if ps -p $$(cat $(PID_FILE)) > /dev/null 2>&1; then \
		echo "Force killing service..."; \
		kill -KILL $$(cat $(PID_FILE)); \
	fi
	@rm -f $(PID_FILE)
	@echo "Service stopped"

# Restart the application
.PHONY: restart
restart:
	@$(MAKE) stop
	@$(MAKE) start

# Clean build artifacts
.PHONY: clean
clean:
	@rm -f $(APP_NAME) $(PID_FILE) $(APP_NAME).log
	@echo "Cleaned build artifacts"

# Show logs
.PHONY: logs
logs:
	@tail -f $(APP_NAME).log
