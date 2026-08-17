package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alextilot/golang-htmx-chatapp/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/database"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
	"github.com/alextilot/golang-htmx-chatapp/internal/router"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
	"github.com/labstack/echo/v5"
)

// main is the only place in the application allowed to read process
// configuration. Everything else receives the specific values it needs via
// constructor parameters (dependency injection), so no other package can
// develop a hidden dependency on the environment.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	db := database.NewDatabase(database.Config{DSN: cfg.DatabasePath})
	defer db.Close()

	// TODO: Move database migrations out of application startup.
	// Production should run migrations through a dedicated migration process.
	db.AutoMigrate()

	authSvc := auth.New(auth.Config{
		JWTSecretKey:        []byte(cfg.JwtSecretKey),
		JWTRefreshSecretKey: []byte(cfg.JwtRefeshSecretKey),
		CookieSecure:        cfg.CookieSecure,
	})

	repos := repository.NewRepositories(repository.Deps{DB: db.Conn})
	services := service.NewServices(service.Deps{Repos: repos})
	handlers := handler.NewHandlers(handler.Deps{Services: services, AuthSvc: authSvc})
	e := router.NewRouter(router.Deps{
		Ctx:                  ctx,
		Handlers:             handlers,
		AuthSvc:              authSvc,
		IsStaticCacheEnabled: cfg.IsProduction(),
	})

	sc := echo.StartConfig{
		Address: cfg.Addr(),
	}

	log.Printf("🌐 Starting server on %s", cfg.Addr())

	if err := sc.Start(ctx, e); err != nil {
		log.Printf("server stopped: %v", err)
	}
}
