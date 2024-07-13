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
	CREATE TABLE IF NOT EXISTS messages (username text not null, content text, time INTEGER);
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

	// e.Use(services.EchoMiddlewareJWTKey())
	e.GET("/", func(etx echo.Context) error {
		user, _ := services.GetUserContext(etx)
		return web.Render(etx, http.StatusOK, views.HomePage(user.IsLoggedIn))
	})

	e.GET("/login", func(etx echo.Context) error {
		user, _ := services.GetUserContext(etx)
		return web.Render(etx, http.StatusOK, views.LoginPage(user.IsLoggedIn))
	})

	e.POST("/login", func(etx echo.Context) error {
		time.Sleep(1 * time.Second)
		return h.Login(etx)
	})

	e.POST("/logout", func(etx echo.Context) error {
		time.Sleep(1 * time.Second)
		return h.Logout(etx)
	})

	e.GET("/signup", func(etx echo.Context) error {
		user, _ := services.GetUserContext(etx)
		return web.Render(etx, http.StatusOK, views.SignUpPage(user.IsLoggedIn))
	})

	e.POST("/signup", func(etx echo.Context) error {
		time.Sleep(1 * time.Second)
		return h.SignUp(etx)
	})

	chatroom := e.Group("/chatroom")
	chatroom.Use(services.EchoMiddlewareJWTConfig())
	chatroom.Use(services.TokenRefresherMiddleware)
	chatroom.GET("", func(etx echo.Context) error {
		user, _ := services.GetUserContext(etx)
		messages, err := messageStore.GetMostRecent(10)
		if err != nil {
			return etx.String(http.StatusBadGateway, "unable to pre populate chat messages")
		}

		var messageViewModel []components.MessageComponentViewModel
		for i := len(messages) - 1; i >= 0; i-- {
			msg := messages[i]
			input := components.MessageComponentViewModel{
				Username: msg.Username,
				Data:     msg.Data,
				Time:     msg.Time.Format("3:04:05 PM"),
				IsSelf:   msg.Username == user.Name,
			}
			messageViewModel = append(messageViewModel, input)
		}
		return web.Render(etx, http.StatusOK, views.ChatroomPage(user.IsLoggedIn, messageViewModel))
	})

	chatroom.GET("/ws", func(etx echo.Context) error {
		user, _ := services.GetUserContext(etx)
		return manager.Handler(etx, ctx, user.Name)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
