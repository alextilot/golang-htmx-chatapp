package router

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/handler/api"

	"github.com/labstack/echo/v5"
)

func registerAPIRoutes(e *echo.Echo, h *api.Handler) {
	v1 := e.Group("/api/v1")

	// Public API routes.

	// Authenticated API routes can be added here using:
	//
	// authenticated := v1.Group("")
	// authenticated.Use(auth.RequireLogin)
	//
	// authenticated.GET("/profile", h.Profile)

	_ = h
	_ = v1
}
