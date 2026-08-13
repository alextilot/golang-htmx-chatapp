package api

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

// Handler contains dependencies required by API handlers.
type Handler struct {
	Services *service.Services
}

// NewHandler creates the API handler with its required dependencies.
func NewHandler(svc *service.Services) *Handler {
	return &Handler{
		Services: svc,
	}
}
