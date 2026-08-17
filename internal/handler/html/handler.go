package html

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

// Handler contains dependencies required by HTML/HTMX handlers.
type Handler struct {
	Services *service.Services
	Auth     *auth.Service
	Group    *GroupHandler
	Chat     *ChatHandler
}

// NewHandler creates the HTML/HTMX handler.
func NewHandler(svc *service.Services, authn *auth.Service) *Handler {
	return &Handler{
		Services: svc,
		Auth:     authn,
		Group:    NewGroupHandler(svc.GroupService),
		Chat:     NewChatHandler(svc.GroupService),
	}
}
