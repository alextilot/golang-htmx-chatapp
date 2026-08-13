package handler

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/handler/api"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler/html"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler/ws"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

// Handlers is the composition container for all HTTP transport handlers.
type Handlers struct {
	HTML *html.Handler
	API  *api.Handler
	WS   *ws.Handler
}

// NewHandlers creates all transport handlers from the application's services.
func NewHandlers(svc *service.Services) *Handlers {
	return &Handlers{
		HTML: html.NewHandler(svc),
		API:  api.NewHandler(svc),
		WS:   ws.NewHandler(svc),
	}
}
