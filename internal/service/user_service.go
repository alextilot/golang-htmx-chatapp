package service

import (
	"context"
	"errors"
	"strings"

	"github.com/alextilot/golang-htmx-chatapp/internal/apperr"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// LoginInput is UserService.Login's own input shape — no transport tags,
// since it must mean the same thing no matter which transport calls it.
type LoginInput struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

// SignupInput is UserService.Signup's own input shape.
type SignupInput struct {
	Username       string `validate:"required,min=2,max=20"`
	Email          string `validate:"required,email"`
	Password       string `validate:"required,password"`
	RepeatPassword string `validate:"required,eqfield=Password"`
}

// GetByID returns a user by ID.
func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := s.repo.FindByID(ctx, id, &user); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperr.NotFound("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Login(ctx context.Context, input LoginInput) (*model.User, error) {
	if appErr := apperr.ValidateStruct(input); appErr != nil {
		return nil, appErr
	}

	user, err := s.repo.GetByUsername(ctx, input.Username)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, apperr.Unauthorized("Invalid login information")
	}
	if err != nil {
		return nil, err
	}
	if !user.CheckPassword(input.Password) {
		return nil, apperr.Unauthorized("Invalid login information")
	}

	return user, nil
}

func (s *UserService) Signup(ctx context.Context, input SignupInput) (*model.User, error) {
	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	fe := apperr.NewFieldErrors()
	if appErr := apperr.ValidateStruct(input); appErr != nil {
		fe = appErr.Fields
	}

	if input.Password != "" && input.Password == input.Username {
		fe.Add("password", "Password cannot match username")
	}

	var reservedUsernames = []string{"admin", "root", "system"}
	for _, word := range reservedUsernames {
		if input.Username != "" && strings.Contains(strings.ToLower(input.Username), word) {
			fe.Add("username", "Username cannot contain reserved words")
		}
	}

	if fe.HasErrors() {
		return nil, apperr.Invalid(fe)
	}

	existing, err := s.repo.FindByUsernameOrEmail(ctx, input.Username, input.Email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		conflict := apperr.NewFieldErrors()
		if existing.Username == input.Username {
			conflict.Add("username", "Username is already taken")
		}
		if existing.Email == input.Email {
			conflict.Add("email", "Email is already taken")
		}
		return nil, apperr.Conflict(conflict)
	}

	return s.createUser(ctx, input.Username, input.Email, input.Password)
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
