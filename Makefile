# Variables allow us to change paths easily later
BINARY_NAME=shiftopt
BUILD_DIR=bin
CMD_PATH=cmd/shiftopt/main.go

# .PHONY tells Make that these are commands, not actual files
.PHONY: all build run clean test help

# Default target: just typing 'make' will run the app
.DEFAULT_GOAL := run

help: ## Show this help message
	@echo 'Usage:'
	@echo '  make build    - Compile the binary to bin/'
	@echo '  make run      - Run the application directly (dev mode)'
	@echo '  make clean    - Remove binary and local database'
	@echo '  make test     - Run unit tests'

build: ## Compile the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build Success! Binary is at $(BUILD_DIR)/$(BINARY_NAME)"

run: build 
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean: ## Clean up build artifacts and reset the DB
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f *.db
	@echo "Clean complete."

test: ## Run all tests in the project
	@go test ./... -v
