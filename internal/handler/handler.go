package handler

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
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

// Deps holds the dependencies needed to construct all transport handlers.
type Deps struct {
	Services *service.Services
	AuthSvc  *auth.Service
}

// NewHandlers creates all transport handlers from the application's
// services and auth.Service.
func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		HTML: html.NewHandler(deps.Services, deps.AuthSvc),
		API:  api.NewHandler(deps.Services),
		WS:   ws.NewHandler(deps.Services),
	}
}
