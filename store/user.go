package store

import (
	"database/sql"
	"log"

	"github.com/alextilot/golang-htmx-chatapp/model"
)

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	store := &UserStore{DB: db}
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (username text not null primary key, password text);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	stmt, err := us.DB.Prepare("SELECT password FROM users WHERE username=?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	var password string
	err = stmt.QueryRow(username).Scan(&password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model.User{
		Username: username,
		Password: password,
	}, nil
}

func (us *UserStore) Create(user *model.User) error {
	sqlStmt := `
	INSERT INTO users(username, password) values(?, ?)
	`
	_, err := us.DB.Exec(sqlStmt, user.Username, user.Password)
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func (us *UserStore) Update(user *model.User) error {
	sqlStmt := `
	UPDATE users SET password=? WHERE username=?
	`
	_, err := us.DB.Exec(sqlStmt, user.Password, user.Username)
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func (us *UserStore) Delete(user *model.User) error {
	sqlStmt := `
	DELETE FROM users WHERE username=?
	`
	_, err := us.DB.Exec(sqlStmt, user.Username)
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return err
	}

	return nil
}
