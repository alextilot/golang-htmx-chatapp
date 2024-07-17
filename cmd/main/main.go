package main

import (
	"context"
	"log"
	"time"

	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/db"
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
	db, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Init services, can log.Fatal
	userStore := store.NewUserStore(db)
	messageStore := store.NewMessageStore(db)

	h := handler.NewHandler(userStore)

	// Init web framework
	e := router.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := NewManager(messageStore)
	go manager.HandleClientListEventChannel(ctx)

	//Init web routes
	e.Static("css", "web/css")
	e.Static("images", "web/images")

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

		var messageViews []*components.MessageComponentViewModel
		for i := len(messages) - 1; i >= 0; i-- {
			msg := messages[i]
			tmp := components.NewMessageView(msg.Username, msg.Data, msg.Time, msg.Username == user.Username)
			messageViews = append(messageViews, tmp)
		}
		return web.Render(etx, http.StatusOK, views.ChatroomPage(user.IsLoggedIn, messageViews))
	})

	chatroom.GET("/ws", func(etx echo.Context) error {
		user, _ := services.GetUserContext(etx)
		return manager.Handler(etx, ctx, user.Username)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
