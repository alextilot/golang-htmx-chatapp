package handler

import "github.com/alextilot/golang-htmx-chatapp/internal/service"

// Handler holds all services that HTTP handlers need.
type Handler struct {
	Services *service.Services
}

// NewHandler creates a new Handler with all required services.
func NewHandler(svc *service.Services) *Handler {
	return &Handler{
		Services: svc,
	}
}
