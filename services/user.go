package services

import (
	"database/sql"
	"golang-app/dto"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *sql.DB
}

func (us *UserService) GetUsers(username string) ([]*dto.UserDto, error) {
	users := []*dto.UserDto{}
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
		users = append(users, &dto.UserDto{
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

func (us *UserService) GetUser(username string) (*dto.UserDto, error) {
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

	return &dto.UserDto{
		ID:       id,
		Username: username,
		Password: password,
	}, nil
}

func (us *UserService) CreateUser(username string, password string) (*dto.UserDto, error) {
	//hash password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return nil, err
	}
	user := &dto.UserDto{
		ID:       uuid.New().String(),
		Username: username,
	}
	sqlStmt := `
	INSERT INTO user(id, username, password) values(?, ?, ?)
	`

	_, err = us.DB.Exec(sqlStmt, user.ID, user.Username, string(hashedPassword))
	if err != nil {
		log.Printf("%q, %s\n", err, sqlStmt)
		return nil, err
	}

	return user, nil
}

func (us *UserService) UpdateUser(id string, username string) (*dto.UserDto, error) {
	user := &dto.UserDto{
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

func (us *UserService) DeleteUser(id string) (*dto.UserDto, error) {
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

func (us *UserService) LoginUser(username string, password string) (*dto.UserDto, error) {
	targetUser, err := us.GetUser(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(password)); err != nil {
		return nil, err
	}

	return targetUser, nil
}
