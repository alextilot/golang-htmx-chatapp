package router

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

// TODO: old remove after implementation.
//	e.GET("/", func(etx echo.Context) error {
//		return web.Render(etx, http.StatusOK, pages.HomePage())
//	})
//
//	e.GET("/login", func(etx echo.Context) error {
//		return web.Render(etx, http.StatusOK, pages.LoginPage())
//	})
//
//	e.GET("/signup", func(etx echo.Context) error {
//		return web.Render(etx, http.StatusOK, pages.SignupPage())
//	})
//
//	e.POST("/login", func(etx echo.Context) error {
//		return handler.Login(etx)
//	})
//
//	e.POST("/signup", func(etx echo.Context) error {
//		return handler.SignUp(etx)
//	})
//
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

func NewRouter(h *handler.Handler) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	// e.Use(middleware.CORS())
	// e.Use(middleware.CSRF())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	// Serve static files
	e.Static("/static", "web/static")

	RegisterPublicRoutes(e, h)
	RegisterPrivateRoutes(e, h)

	return e
}
