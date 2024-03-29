# run:
# 	@templ generate
# 	@npx tailwind -i 'css/styles.css' -o 'css/tailwind.css'
# 	@go run main.go

.PHONY: templ
templ:
	templ generate -watch -proxy=http://localhost:3000

.PHONY: tailwind
tailwind:
	npx tailwindcss -i ./css/input.css -o ./css/output.css --watch

.PHONY: air
air: 
	air -c ./.air.toml

# install:
#   @go install github.com/a-h/templ/cmd/templ@latest
# 	@go get ./...
# 	@go mod vendor
# 	@go mod tidy
# 	@go mod download 
# 	@npm i

# build:
# 	tailwindcss -i css/main.css -o css/styles.css
# 	@templ generate view
# 	@go build -o bin/github.com/alextilot/golang-htmx-chatapp main.go