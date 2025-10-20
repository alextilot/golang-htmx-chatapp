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
