package router

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterPrivateRoutes(e *echo.Echo, h *handler.Handler, ctx context.Context) {
	private := e.Group("/app")
	// TODO: re-add your auth middleware here once ready.
	// private.Use(auth.RequireLogin)

	// ---- Groups ----
	private.GET("/groups", h.Group.List)
	private.POST("/groups", h.Group.Create)
	private.GET("/groups/:groupID", h.Group.Show)
	private.PUT("/groups/:groupID", h.Group.Update)
	private.DELETE("/groups/:groupID", h.Group.Delete)
	private.POST("/groups/:groupID/members", h.Group.AddMember)

	// ---- WebSocket ----
	// Sits outside /app so the WS upgrade is not affected by session middleware.
	e.GET("/ws/groups/:groupID", func(c echo.Context) error {
		return h.Group.ConnectWS(c, ctx)
	})

}

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

// guardedRoutes := e.Group("/chatroom")
// guardedRoutes.Use(services.TokenRefresherMiddleware)
//
//	guardedRoutes.Use(echojwt.WithConfig(echojwt.Config{
//		SigningKey:   []byte(config.Cfg.JwtSecretKey),
//		TokenLookup:  "cookie:access-token",
//		ErrorHandler: services.JWTErrorChecker,
//	}))
//
//	guardedRoutes.GET("", func(etx echo.Context) error {
//		return web.Render(etx, http.StatusOK, pages.ChatroomPage())
//	})
//
//	e.GET("/ws/chatroom", func(etx echo.Context) error {
//		return manager.Handler(etx, ctx)
//	})

// ctx, cancel := context.WithCancel(context.Background())
// defer cancel()
//
// manager := NewManager()
// go manager.HandleClientListEventChannel(ctx)
