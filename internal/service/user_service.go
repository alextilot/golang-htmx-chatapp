package service

import (
	"errors"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Login handles login logic
func (s *UserService) Login(username string, password string) (*model.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

// Creates a user and hashes the password
func (s *UserService) Create(username, email, password string) (*model.User, error) {
	user := &model.User{
		Username: username,
		Email:    email,
	}

	hashed, err := user.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Password = hashed

	// Persist user via repository
	return s.repo.Create(user)
}

func (s *UserService) Signup(username, email, password string) (*model.User, validation.FieldErrors, error) {
	fe := validation.FieldErrors{}

	// 1. Check for existing username/email
	existing, err := s.repo.ExistsByUsernameOrEmail(username, email)
	if err != nil {
		return nil, nil, err
	}
	if existing != nil {
		if existing.Username == username {
			fe["username"] = []string{"Username is already taken"}
		}
		if existing.Email == email {
			fe["email"] = []string{"Email is already taken"}
		}
		return nil, fe, nil
	}

	// 2. Delegate creation to Create() method (handles password hashing)
	user, err := s.Create(username, email, password)
	if err != nil {
		return nil, nil, err
	}

	return user, nil, nil
}
