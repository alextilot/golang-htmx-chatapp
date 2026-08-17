package ws

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

// Handler contains dependencies required by WebSocket handlers.
type Handler struct {
	Services *service.Services
	Hub      *Hub
	Group    *GroupHandler
	Chat     *ChatHandler
}

// NewHandler creates the WebSocket handler and its connection hub.
func NewHandler(svc *service.Services) *Handler {
	hub := NewHub(svc.GroupService)

	return &Handler{
		Services: svc,
		Hub:      hub,
		Group:    NewGroupHandler(hub, svc.GroupService),
		Chat:     NewChatHandler(hub, svc.GroupService),
	}
}
