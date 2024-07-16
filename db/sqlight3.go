package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const dbSourceName string = "./db/file.db"

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
