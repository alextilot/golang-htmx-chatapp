package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const dbSourceName string = "./main.db"

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbSourceName)
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (username text not null primary key, password text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	sqlStmt = `
	CREATE TABLE IF NOT EXISTS messages (username text not null, content text, time INTEGER);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}
