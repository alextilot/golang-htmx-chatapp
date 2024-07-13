package handler

import "github.com/alextilot/golang-htmx-chatapp/store"

type Handler struct {
	userService *store.UserStore
}

func NewHandler(us *store.UserStore) *Handler {
	return &Handler{
		userService: us,
	}
}
