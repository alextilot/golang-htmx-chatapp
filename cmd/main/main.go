package main

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/database"
	"github.com/alextilot/golang-htmx-chatapp/internal/handler"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
	"github.com/alextilot/golang-htmx-chatapp/internal/router"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
	"log"
)

func main() {
	// Connect to the database
	db := database.New(config.Cfg.DatabasePath)
	defer db.Close()

	// Run migrations for all registered models
	db.AutoMigrate()

	r := repository.NewRepositories(db.Conn)
	s := service.NewServices(r)
	h := handler.NewHandler(s)
	e := router.NewRouter(h)

	// Start the web server
	log.Printf("🚀 Starting server on port %s", config.Cfg.Port)
	e.Logger.Fatal(e.Start(":" + config.Cfg.Port))
}
