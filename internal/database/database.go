package database

import (
	"log"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Conn *gorm.DB
}

// Config holds the settings needed to open a database connection.
type Config struct {
	DSN string
}

func NewDatabase(cfg Config) *Database {
	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &Database{Conn: db}
}

func (d *Database) AutoMigrate() {
	if err := d.Conn.AutoMigrate(model.Models()...); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func (d *Database) Close() error {
	sqlDB, err := d.Conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
