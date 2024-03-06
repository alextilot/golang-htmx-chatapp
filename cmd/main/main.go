package main

import (
	"context"
	"database/sql"
	"log"

	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/handler"
	"github.com/alextilot/golang-htmx-chatapp/router"
	"github.com/alextilot/golang-htmx-chatapp/services"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/views"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Init DB
	db, err := sql.Open("sqlite3", "./db/main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS user (id text not null primary key, username text, password text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s \n", err, sqlStmt)
		return
	}
	// Init services
	userService := &services.UserService{
		DB: db,
	}

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
		return handler.PostLogin(etx, ctx, userService)
	})
	e.POST("/signup", func(etx echo.Context) error {
		return handler.PostSignup(etx, ctx, userService)
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
