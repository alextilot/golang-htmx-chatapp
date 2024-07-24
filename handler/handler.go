package handler

import "github.com/alextilot/golang-htmx-chatapp/store"

type Handler struct {
	userService *store.UserStore
	msgService  *store.MessageStore
}

func NewHandler(us *store.UserStore, ms *store.MessageStore) *Handler {
	return &Handler{
		userService: us,
		msgService:  ms,
	}
}
