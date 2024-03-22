package handler

import "github.com/alextilot/golang-htmx-chatapp/services"

type Handler struct {
	userService *services.UserService
}

func NewHandler(us *services.UserService) *Handler {
	return &Handler{
		userService: us,
	}
}
