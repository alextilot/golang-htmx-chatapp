package router

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterPrivateRoutes(e *echo.Echo, h *handler.Handler) {
	// Protected routes (JWT middleware)
	// private := e.Group("/app")
	// private.Use(JWTMiddleware()) // implement JWT middleware
	//
	// private.GET("/chatroom", h.ChatRoom)
	// private.POST("/chat/message", h.SendMessage)
	// private.GET("/chat/messages", h.GetMessages)
	//
	// // Users
	// private.GET("/users", h.ListUsers)
	// private.GET("/users/:id", h.GetUser)
}
