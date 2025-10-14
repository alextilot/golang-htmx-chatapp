package database

import (
	"github.com/alextilot/golang-htmx-chatapp/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dsn = "./app_main.sqlight3"
)

// DB holds the GORM database instance
var DB *gorm.DB

// Models for the database
var models = []any{
	&model.User{},
	&model.ChatGroup{},
	&model.ChatMessage{},
}

// Connect initializes the database connection
func Connect() error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

// Migrate performs automatic database migrations
func Migrate() error {
	if err := DB.AutoMigrate(models...); err != nil {
		return err
	}
	return nil
}

// Close closes the database connection
func Close() error {
	db, err := DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
