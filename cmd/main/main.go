package main

import (
	"context"
	"log"

	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/database"
	"github.com/alextilot/golang-htmx-chatapp/handler"
	"github.com/alextilot/golang-htmx-chatapp/router"
	"github.com/alextilot/golang-htmx-chatapp/services"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/views"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal(err)
	}

	defer database.Close()

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
		return web.Render(etx, http.StatusOK, views.HomePage())
	})
	e.GET("/login", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, views.LoginPage())
	})
	e.GET("/signup", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, views.SignupPage())
	})

	e.POST("/login", func(etx echo.Context) error {
		return handler.Login(etx, ctx)
	})
	e.POST("/signup", func(etx echo.Context) error {
		return handler.SignUp(etx, ctx)
	})

	guardedRoutes := e.Group("/chatroom")
	guardedRoutes.Use(services.TokenRefresherMiddleware)
	guardedRoutes.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:   []byte(services.JwtSecretKey),
		TokenLookup:  "cookie:access-token",
		ErrorHandler: services.JWTErrorChecker,
	}))
	guardedRoutes.GET("", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, views.ChatroomPage())
	})

	e.GET("/ws/chatroom", func(etx echo.Context) error {
		return manager.Handler(etx, ctx)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
