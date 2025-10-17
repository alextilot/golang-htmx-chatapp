package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Conn *gorm.DB
}

func New(dns string) *Database {
	db, err := gorm.Open(sqlite.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &Database{Conn: db}
}

func (d *Database) Migrate(models ...any) {
	if err := d.Conn.AutoMigrate(models...); err != nil {
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
