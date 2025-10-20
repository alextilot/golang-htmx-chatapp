# Makefile for Go + HTMX + Tailwind v4 + Templ

# Run the full dev environment: templ watcher + tailwind watcher + air
.PHONY: dev
dev:
	@echo "Starting dev environment..."
	$(MAKE) tailwind &
	$(MAKE) air

# Generate templ files in watch mode
.PHONY: templ
templ:
	go tool templ generate --watch --proxy=http://localhost:8080

# Watch Tailwind v4 CSS
.PHONY: tailwind
tailwind:
	npm run watch

# Run Go server with air (live reload)
.PHONY: air
air: 
	go tool air -c ./.air.toml

# Build production Tailwind CSS
.PHONY: css
css:
	npm run build

# Clean temp directories and generated files
.PHONY: clean
clean:
	rm -rf ./tmp ./web/static/*.dist.*

