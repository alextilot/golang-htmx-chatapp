package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/handler"
	"github.com/alextilot/golang-htmx-chatapp/router"
	"github.com/alextilot/golang-htmx-chatapp/services"
	"github.com/alextilot/golang-htmx-chatapp/store"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/components"
	"github.com/alextilot/golang-htmx-chatapp/web/views"

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
	CREATE TABLE IF NOT EXISTS user (username text not null primary key, password text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s \n", err, sqlStmt)
		return
	}

	sqlStmt = `
	CREATE TABLE IF NOT EXISTS messages (clientId text not null, username text not null, content text, time INTEGER);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s \n", err, sqlStmt)
		return
	}

	// Init services
	userStore := &store.UserStore{
		DB: db,
	}

	messageStore := &store.MessageStore{
		DB: db,
	}

	h := handler.NewHandler(userStore)

	// Init web framework
	e := router.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := NewManager(messageStore)
	go manager.HandleClientListEventChannel(ctx)

	//Init web routes
	e.Static("/css", "css")

	e.GET("/", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, views.HomePage())
	})

	e.GET("/login", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, views.LoginPage())
	})

	e.POST("/login", func(etx echo.Context) error {
		time.Sleep(1 * time.Second)
		return h.Login(etx)
	})

	e.GET("/signup", func(etx echo.Context) error {
		return web.Render(etx, http.StatusOK, views.SignUpPage())
	})

	e.POST("/signup", func(etx echo.Context) error {
		time.Sleep(1 * time.Second)
		return h.SignUp(etx)
	})

	chatroom := e.Group("/chatroom")
	chatroom.Use(services.TokenRefresherMiddleware)
	chatroom.Use(services.EchoMiddlewareJWTConfig())
	chatroom.GET("", func(etx echo.Context) error {
		name := services.GetUsername(etx)
		messages, err := messageStore.GetLast(10)
		if err != nil {
			return etx.String(http.StatusBadGateway, "unable to pre populate chat messages")
		}

		var messageViewModel []components.MessageComponentViewModel
		for _, msg := range messages {
			input := components.MessageComponentViewModel{
				Username: msg.Username,
				Data:     msg.Data,
				Time:     msg.Time.Format("3:04:05 PM"),
				IsSelf:   msg.Username == name,
			}
			messageViewModel = append(messageViewModel, input)
		}
		return web.Render(etx, http.StatusOK, views.ChatroomPage(messageViewModel))
	})
	chatroom.GET("/ws", func(etx echo.Context) error {
		name := services.GetUsername(etx)
		return manager.Handler(etx, ctx, name)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
