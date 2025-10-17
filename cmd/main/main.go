package main

import (
	"context"

	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/config"
	"github.com/alextilot/golang-htmx-chatapp/database"
	"github.com/alextilot/golang-htmx-chatapp/handler"
	"github.com/alextilot/golang-htmx-chatapp/model"
	"github.com/alextilot/golang-htmx-chatapp/router"
	"github.com/alextilot/golang-htmx-chatapp/services"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/pages"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

var (
	models = []any{
		&model.User{},
		&model.ChatGroup{},
		&model.ChatMessage{},
	}
)

func main() {
	// Connect database
	db := database.New(config.Cfg.DatabasePath)
	defer db.Close()

	// Run migrations
	db.Migrate(models...)

	// Init services
	userService := &services.UserService{
		DB: database.DB,
	}

	handler := handler.NewHandler(userService)

	// Init web framework
	e := router.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := NewManager()
	go manager.HandleClientListEventChannel(ctx)

	//Init web routes
	e.Static("/css", "css")

	e.GET("/", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, pages.HomePage())
	})
	e.GET("/login", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, pages.LoginPage())
	})
	e.GET("/signup", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, pages.SignupPage())
	})

	e.POST("/login", func(etx echo.Context) error {
		return handler.Login(etx)
	})
	e.POST("/signup", func(etx echo.Context) error {
		return handler.SignUp(etx)
	})

	guardedRoutes := e.Group("/chatroom")
	guardedRoutes.Use(services.TokenRefresherMiddleware)
	guardedRoutes.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:   []byte(config.Cfg.JwtSecretKey),
		TokenLookup:  "cookie:access-token",
		ErrorHandler: services.JWTErrorChecker,
	}))
	guardedRoutes.GET("", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, pages.ChatroomPage())
	})

	e.GET("/ws/chatroom", func(etx echo.Context) error {
		return manager.Handler(etx, ctx)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
