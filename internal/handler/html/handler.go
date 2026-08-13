package html

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

// Handler contains dependencies required by HTML/HTMX handlers.
type Handler struct {
	Services *service.Services
	Group    *GroupHandler
}

// NewHandler creates the HTML/HTMX handler.
func NewHandler(svc *service.Services) *Handler {
	return &Handler{
		Services: svc,
		Group:    NewGroupHandler(svc.GroupService),
	}
}
