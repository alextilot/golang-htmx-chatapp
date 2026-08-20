# Makefile for Go + HTMX + Tailwind v4 + Templ + Air

# Detect OS-specific Air config
ifeq ($(OS),Windows_NT)
	AIR_CONFIG := ./.air.windows.toml
else
	AIR_CONFIG := ./.air.unix.toml
endif

# Default seed target: a scratch DB, so seeding never touches real dev data.
# Override to point at another file, e.g.:
#   make db-seed DATABASE_PATH=./internal/db/app.main.sqlite3
DATABASE_PATH ?= ./internal/db/seed.sample.sqlite3

# ---------------------------------------------
# One-time setup: install Go tools, tidy modules, install npm deps
# ---------------------------------------------
.PHONY: setup
setup:
	@echo "📦 Installing dependencies..."
	go mod tidy
	npm install
	@echo "✅ Setup complete. Run 'make dev' to start."

# ---------------------------------------------
# Run full dev environment: templ watcher + tailwind watcher + air
# Ctrl+C stops all three (trap kills the whole process group on exit)
# ---------------------------------------------
.PHONY: dev
dev:
	@echo "🚀 Starting dev environment..."
	@trap 'kill 0' EXIT INT TERM; \
	$(MAKE) templ & \
	$(MAKE) tailwind & \
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
# Build minified production CSS
# ---------------------------------------------
.PHONY: css
css:
	@echo "🎨 Building production CSS..."
	npm run build

# ---------------------------------------------
# Run Go server with Air (live reload)
# Use AIR_CONFIG env var to switch config
# ---------------------------------------------
.PHONY: air
air: 
	@echo "💨 Starting Air with config: $(AIR_CONFIG)"
	go tool air -c $(AIR_CONFIG)


#  ---------------------------------------------
# Drop and reseed the database from cmd/seed/data/*.json (see DATABASE_PATH
# above) — always starts clean, so it's safe to rerun after manual testing.
# ---------------------------------------------
.PHONY: db-seed
db-seed:
	@echo "🌱 Dropping and reseeding $(DATABASE_PATH)..."
	@rm -f $(DATABASE_PATH)
	DATABASE_PATH=$(DATABASE_PATH) go run ./cmd/seed

#  ---------------------------------------------
# Clean temp directories and generated files
# ---------------------------------------------
.PHONY: clean
clean:
	@echo "🧼 Cleaning temporary and generated files..."
	@rm -rf ./tmp
	@rm -f ./web/static/*.dist.*
	@find ./web -type f -name "*_templ.go" -exec rm -f {} +
	@echo "✅ Clean complete."
