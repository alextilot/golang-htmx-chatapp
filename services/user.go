package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/alextilot/golang-htmx-chatapp/model"

	"github.com/google/uuid"
)

type UserService struct {
	DB *sql.DB
}

func (us *UserService) GetUsers(username string) ([]*model.User, error) {
	users := []*model.User{}
	rows, err := us.DB.Query("SELECT id, username FROM user WHERE username = ?", username)
	if err != nil {
		log.Println(err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var username string
		err = rows.Scan(&id, &username)
		if err != nil {
			log.Println(err)
			return users, err
		}
		users = append(users, &model.User{
			ID:       id,
			Username: username,
		})
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return users, err
	}
	return users, nil
}

func (us *UserService) GetUser(username string) (*model.User, error) {
	stmt, err := us.DB.Prepare("SELECT id, password FROM user WHERE username=?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	var id string
	var password string
	err = stmt.QueryRow(username).Scan(&id, &password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model.User{
		ID:       id,
		Username: username,
		Password: password,
	}, nil
}

func (us *UserService) CreateUser(username string, password string) (*model.User, error) {
	user := &model.User{
		ID:       uuid.New().String(),
		Username: username,
	}
	//hash password.
	hashedPassword, err := user.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	sqlStmt := `
	INSERT INTO user(id, username, password) values(?, ?, ?)
	`

	_, err = us.DB.Exec(sqlStmt, user.ID, user.Username, user.Password)
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return nil, err
	}

	return user, nil
}

func (us *UserService) UpdateUser(id string, username string) (*model.User, error) {
	user := &model.User{
		ID:       id,
		Username: username,
	}
	sqlStmt := `
	UPDATE user SET username=? WHERE id=?
	`

	_, err := us.DB.Exec(sqlStmt, user.Username, user.ID)
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return nil, err
	}

	return user, nil
}

func (us *UserService) DeleteUser(id string) (*model.User, error) {
	sqlStmt := `
	DELETE FROM user WHERE id=?
	`

	_, err := us.DB.Exec(sqlStmt, id)
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return nil, err
	}

	return nil, nil
}

func (us *UserService) LoginUser(username string, password string) (*model.User, error) {
	user, err := us.GetUser(username)
	if err != nil {
		return nil, err
	}

	if user.CheckPassword(password) {
		return user, nil
	}

	return nil, errors.New("invalid login information")
}
