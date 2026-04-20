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
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := config.Load()

	db := database.New(config.Cfg.DatabasePath)
	defer db.Close()
	db.AutoMigrate()

	r := repository.NewRepositories(db.Conn)
	s := service.NewServices(r)
	h := handler.NewHandler(s)
	e := router.NewRouter(h, ctx)

	// Shut down the Echo server cleanly when ctx is cancelled (Ctrl+C / SIGTERM).
	go func() {
		<-ctx.Done()
		log.Println("shutting down...")
		if err := e.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Printf("🌐 Starting server on %s", cfg.Addr())
	if err := e.Start(cfg.Addr()); err != nil {
		log.Printf("server stopped: %v", err)
	}
}
