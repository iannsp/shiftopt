# Variables allow us to change paths easily later
BINARY_NAME=shiftopt
BUILD_DIR=bin
CMD_PATH=cmd/shiftopt/main.go

# .PHONY tells Make that these are commands, not actual files
.PHONY: all build clean test help

# Default: Build Both


all: build

help: ## Show this help message
	@echo 'Usage:'
	@echo '  make build    - Compile the binary to bin/'
	@echo '  make run      - Run the application directly (dev mode)'
	@echo '  make clean    - Remove binary and local database'
	@echo '  make test     - Run unit tests'


build:
	@echo "Building ShiftOpt (CSV Generator)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/shiftopt cmd/shiftopt/main.go
	
	@echo "Building ShiftSummary (Diagnostic Tool)..."
	@go build -o $(BUILD_DIR)/shiftsummary cmd/shiftsummary/main.go
	
	@echo "Build Complete. Artifacts in $(BUILD_DIR)/"

# Run the summary by default
run: build
	@./$(BUILD_DIR)/shiftsummary

# Run the export
export: build
	@./$(BUILD_DIR)/shiftopt

clean:
	@rm -rf $(BUILD_DIR)
	@rm -f *.db
	@rm -f *.csv

test:
	@go test ./... -v
