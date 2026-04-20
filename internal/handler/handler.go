package handler

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/handler/ws"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

// Handler holds all services and sub-handlers that HTTP handlers need.
type Handler struct {
	Services *service.Services
	Hub      *ws.Hub
	Group    *GroupHandler
}

// NewHandler creates a new Handler with all required services.
func NewHandler(svc *service.Services) *Handler {
	hub := ws.NewHub()

	return &Handler{
		Services: svc,
		Hub:      hub,
		Group:    NewGroupHandler(hub, svc.GroupService),
	}
}
