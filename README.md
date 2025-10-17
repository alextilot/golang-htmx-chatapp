# HTMX Go Prototype

Server-side UI prototype built with:

- [Go](https://golang.org/)
- [HTMX](https://htmx.org/)
- [Templ](https://github.com/a-h/templ)
- [Tailwind CSS v4](https://tailwindcss.com/)
- [daisyUI](https://daisyui.com/)
- [Air](https://github.com/cosmtrek/air) for live reload

## Project Structure
<!-- TODO:  finish this section -->

## Development

Install dependencies:

```bash
go install github.com/a-h/templ/cmd/templ@latest
go mod tidy
npm install
```

Start the development environment:

```bash
make dev
```

- templ → generates templates in watch mode

- tailwind → watches CSS changes

- air → live reloads the Go server

Build production CSS:
```bash
make css
```

Clean temporary files:
```bash
make clean
```

Notes

Tailwind v4 is CSS-first: all configuration lives in app.css.

DaisyUI is included for pre-built components and themes.

air watches Go and template files for instant feedback.

