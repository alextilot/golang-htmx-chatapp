ifeq ($(OS),Windows_NT)
    uname_S := Windows
else
    uname_S := $(shell uname -s)
endif

ifeq ($(uname_S), Windows)
    air_tomel = ./.air.toml_Windows
endif
ifeq ($(uname_S), Linux)
    air_tomel = ./.air.toml_Unix 
endif

.PHONY: templ
templ:
	templ generate --watch

.PHONY: tailwind
tailwind:
	npx tailwindcss -i ./web/css/global.css -o ./web/css/dist.css --watch

.PHONY: air
air: 
	air -c ${air_tomel}

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

# run:
# 	@templ generate
# 	@npx tailwind -i 'css/styles.css' -o 'css/tailwind.css'
# 	@go run main.go
