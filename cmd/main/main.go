package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/database"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
	"github.com/alextilot/golang-htmx-chatapp/internal/router"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
	"github.com/labstack/echo/v5"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	db := database.New(config.Cfg.DatabasePath)
	defer db.Close()

	// TODO: Move database migrations out of application startup.
	// Production should run migrations through a dedicated migration process.
	db.AutoMigrate()

	repos := repository.NewRepositories(db.Conn)
	services := service.NewServices(repos)
	handlers := handler.NewHandlers(services)
	e := router.NewRouter(ctx, handlers)

	sc := echo.StartConfig{
		Address: cfg.Addr(),
	}

	log.Printf("🌐 Starting server on %s", cfg.Addr())

	if err := sc.Start(ctx, e); err != nil {
		log.Printf("server stopped: %v", err)
	}
}
