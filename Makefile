# Makefile for Go + HTMX + Tailwind v4 + Templ + Air
#
# Run full dev environment: templ watcher + tailwind watcher + air
.PHONY: dev
dev:
	@echo "🚀 Starting dev environment..."
	$(MAKE) templ &
	$(MAKE) tailwind &
	$(MAKE) air

# Generate templ files in watch mode
.PHONY: templ
templ:
	@echo "👀 Watching templ files..."
	go tool templ generate --watch --open-browser=false

# Watch Tailwind v4 CSS
.PHONY: tailwind
tailwind:
	@echo "🎨 Watching Tailwind..."
	npm run watch

# Run Go server with air (live reload)
.PHONY: air
air: 
	@echo "💨 Starting Air..."
	go tool air -c ./.air.toml

# Build production Tailwind CSS
.PHONY: css
css:
	npm run build

# Clean temp directories and generated files
.PHONY: clean
clean:
	rm -rf ./tmp ./web/static/*.dist.*

