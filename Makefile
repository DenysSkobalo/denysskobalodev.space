# Project Configuration & Variables
BINARY_NAME  := main.wasm
BUILD_DIR    := build
SRC_DIR      := src
CMD_DIR      := ./cmd/worker/...
GO_VERSION   := release-branch.go1.22
WASM_EXEC_URL:= https://raw.githubusercontent.com/golang/go/$(GO_VERSION)/misc/wasm/wasm_exec.js

# Terminal Colors for Output
CYAN   := \033[0;36m
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
RESET  := \033[0m

.PHONY: all help dev build fetch-runtime clean test check-deps d1-migrate-local tail-backend tail-config

all: help

## help: Display available targets
help:
	@echo "$(CYAN)Usage: make [target]$(RESET)"
	@echo ""
	@echo "$(GREEN)Targets:$(RESET)"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## dev: Full build pipeline and launch local Wrangler Edge emulator
dev: check-deps fetch-runtime build
	@echo "$(CYAN)[INFO] Starting local Cloudflare Workers emulator via Wrangler...$(RESET)"
	npx wrangler dev

## build: Compile Go code to WebAssembly with size optimizations (-s -w)
build:
	@echo "$(CYAN)[INFO] Creating build directory structure...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@echo "$(CYAN)[INFO] Compiling Go WebAssembly binary for GOOS=js GOARCH=wasm...$(RESET)"
	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)[SUCCESS] WebAssembly binary successfully compiled: $(BUILD_DIR)/$(BINARY_NAME)$(RESET)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME) | awk '{print "$(YELLOW)[SIZE] WASM Payload Size: " $$5 "$(RESET)"}'

## fetch-runtime: Download compatible Go JS WASM bridge runtime (wasm_exec.js)
fetch-runtime:
	@echo "$(CYAN)[INFO] Ensuring JS source directory exists...$(RESET)"
	@mkdir -p $(SRC_DIR)
	@if [ ! -f $(SRC_DIR)/wasm_exec.js ]; then \
		echo "$(YELLOW)[WARN] wasm_exec.js missing. Fetching Go WASM bridge from upstream...$(RESET)"; \
		curl -sSL "$(WASM_EXEC_URL)" -o $(SRC_DIR)/wasm_exec.js; \
		echo "$(GREEN)[SUCCESS] Downloaded wasm_exec.js to $(SRC_DIR)/$(RESET)"; \
	else \
		echo "$(GREEN)[OK] Go WASM runtime bridge $(SRC_DIR)/wasm_exec.js already exists.$(RESET)"; \
	fi

## d1-migrate-local: Apply Cloudflare D1 database migrations locally
d1-migrate-local:
	@echo "$(CYAN)[INFO] Applying local D1 database migrations...$(RESET)"
	npx wrangler d1 migrations apply DB --local -c wrangler.toml

## tail-backend: Stream real-time Worker logs by service name
tail-backend:
	@echo "$(CYAN)[INFO] Tailing logs for denysskobalo-backend-worker...$(RESET)"
	npx wrangler tail denysskobalo-backend-worker

## tail-config: Stream real-time logs based on wrangler.toml configuration
tail-config:
	@echo "$(CYAN)[INFO] Tailing logs using wrangler.toml config...$(RESET)"
	npx wrangler tail -c wrangler.toml

## test: Execute Go unit and integration tests with race detection
test:
	@echo "$(CYAN)[INFO] Running unit and integration tests...$(RESET)"
	go test -v -race ./tests/...

## check-deps: Verify toolchain dependencies (Go, Node.js, Wrangler)
check-deps:
	@echo "$(CYAN)[INFO] Checking toolchain dependencies...$(RESET)"
	@command -v go >/dev/null 2>&1 || { echo "$(RED)[ERROR] Go is not installed. Aborting.$(RESET)"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "$(RED)[ERROR] Node.js is not installed. Aborting.$(RESET)"; exit 1; }
	@command -v npx >/dev/null 2>&1 || { echo "$(RED)[ERROR] npx/npm is not installed. Aborting.$(RESET)"; exit 1; }
	@echo "$(GREEN)[OK] Toolchain requirements verified.$(RESET)"

## clean: Clean up build artifacts, local Wrangler state, and caches
clean:
	@echo "$(YELLOW)[INFO] Cleaning build outputs and local caches...$(RESET)"
	@rm -rf $(BUILD_DIR)
	@rm -rf .wrangler
	@rm all_files.txt
	@echo "$(GREEN)[SUCCESS] Cleanup completed.$(RESET)"
