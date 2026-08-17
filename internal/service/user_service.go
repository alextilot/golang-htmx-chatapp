package service

import (
	"context"
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

// GetByID returns a user by ID. Returns an error if not found.
func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := s.repo.FindByID(ctx, id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil || !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) Signup(
	ctx context.Context,
	username, email, password string,
) (*model.User, validation.FieldErrors, error) {

	fe := validation.FieldErrors{}

	existing, err := s.repo.FindByUsernameOrEmail(ctx, username, email)
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

	user, err := s.createUser(ctx, username, email, password)
	if err != nil {
		return nil, nil, err
	}

	return user, nil, nil
}

func (s *UserService) createUser(ctx context.Context, username, email, password string) (*model.User, error) {
	user := &model.User{
		Username: username,
		Email:    email,
	}

	hashed, err := user.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user.Password = hashed

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
