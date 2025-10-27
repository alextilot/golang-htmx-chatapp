# Makefile for Go + HTMX + Tailwind v4 + Templ + Air

# Default Air config
AIR_CONFIG ?= ./.air.toml
# Detect OS
ifeq ($(OS),Windows_NT)
	AIR_CONFIG := ./.air.windows.toml
else
	AIR_CONFIG := ./.air.unix.toml
endif

# ---------------------------------------------
# Run full dev environment: templ watcher + tailwind watcher + air
# ---------------------------------------------
.PHONY: dev
dev:
	@echo "🚀 Starting dev environment..."
	$(MAKE) templ &
	$(MAKE) tailwind &
	$(MAKE) air

# ---------------------------------------------
# Generate templ files in watch mode
# ---------------------------------------------
.PHONY: templ
templ:
	@echo "👀 Watching templ files..."
	go tool templ generate --watch --open-browser=false

# ---------------------------------------------
# Watch Tailwind v4 CSS
# ---------------------------------------------
.PHONY: tailwind
tailwind:
	@echo "🎨 Watching Tailwind..."
	npm run watch

# ---------------------------------------------
# Run Go server with Air (live reload)
# Use AIR_CONFIG env var to switch config
# ---------------------------------------------
.PHONY: air
air: 
	@echo "💨 Starting Air with config: $(AIR_CONFIG)"
	go tool air -c $(AIR_CONFIG)


#  ---------------------------------------------
# Clean temp directories and generated files
# ---------------------------------------------
.PHONY: clean
clean:
	@echo "🧼 Cleaning temporary and generated files..."
	@rm -rf ./tmp
	@rm -f ./web/static/*.dist.*
	@find ./web -type f -name "*.templ.go" -exec rm -f {} +
	@echo "✅ Clean complete."
