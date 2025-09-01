# Makefile for AstroEph API
# Astrological Calculation Service

# Variables
BINARY_NAME=astroeph-api
MAIN_PACKAGE=.
PORT=8080

# Swiss Ephemeris library path (required for runtime)
export DYLD_LIBRARY_PATH=/usr/local/lib

# Default target
.PHONY: help
help:
	@echo "AstroEph API - Available Commands:"
	@echo ""
	@echo "  make run        - Start the development server"
	@echo "  make build      - Build the binary"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make deps       - Download and tidy dependencies"
	@echo "  make health     - Check if server is running"
	@echo "  make natal      - Test natal chart endpoint"
	@echo ""

# Run the development server with Swiss Ephemeris support
.PHONY: run
run:
	@echo "🌟 Starting AstroEph API server..."
	@echo "📍 Port: $(PORT)"
	@echo "📚 Library Path: $(DYLD_LIBRARY_PATH)"
	@echo "🔗 Health Check: http://localhost:$(PORT)/health"
	@echo "📖 API Endpoint: http://localhost:$(PORT)/api/v1/natal-chart"
	@echo ""
	DYLD_LIBRARY_PATH=$(DYLD_LIBRARY_PATH) go run $(MAIN_PACKAGE)

# Build the binary
.PHONY: build
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "✅ Built: ./$(BINARY_NAME)"

# Run the built binary
.PHONY: run-binary
run-binary: build
	@echo "🚀 Running built binary..."
	DYLD_LIBRARY_PATH=$(DYLD_LIBRARY_PATH) ./$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	go clean
	rm -f $(BINARY_NAME)
	@echo "✅ Clean complete"

# Download and tidy dependencies
.PHONY: deps
deps:
	@echo "📦 Managing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies updated"

# Run tests
.PHONY: test
test:
	@echo "🧪 Running tests..."
	DYLD_LIBRARY_PATH=$(DYLD_LIBRARY_PATH) go test -v ./...

# Health check
.PHONY: health
health:
	@echo "🏥 Checking server health..."
	@curl -s http://localhost:$(PORT)/health | jq . || echo "❌ Server not responding"

# Test natal chart endpoint
.PHONY: natal
natal:
	@echo "🌟 Testing natal chart endpoint..."
	@curl -X POST http://localhost:$(PORT)/api/v1/natal-chart \
		-H "Content-Type: application/json" \
		-d '{ \
			"day": 1, \
			"month": 1, \
			"year": 2000, \
			"local_time": "12:00:00", \
			"city": "London", \
			"house_system": "Placidus" \
		}' | jq . || echo "❌ Natal chart endpoint not responding"

# Development setup check
.PHONY: check-deps
check-deps:
	@echo "🔍 Checking Swiss Ephemeris installation..."
	@if [ -f "/usr/local/lib/libswe.dylib" ]; then \
		echo "✅ Swiss Ephemeris library found"; \
	else \
		echo "❌ Swiss Ephemeris library not found at /usr/local/lib/libswe.dylib"; \
		echo "   Please run the Swiss Ephemeris setup first."; \
		exit 1; \
	fi
	@echo "🔍 Checking Go dependencies..."
	@go mod verify
	@echo "✅ All dependencies verified"

# Install (for production deployment)
.PHONY: install
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"

# Uninstall
.PHONY: uninstall
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstalled"

# Development mode with auto-restart (requires 'air' tool)
.PHONY: dev
dev:
	@if command -v air > /dev/null; then \
		echo "🔄 Starting development mode with auto-restart..."; \
		DYLD_LIBRARY_PATH=$(DYLD_LIBRARY_PATH) air; \
	else \
		echo "📝 Install 'air' for auto-restart: go install github.com/cosmtrek/air@latest"; \
		echo "🔄 Using standard run mode..."; \
		$(MAKE) run; \
	fi
