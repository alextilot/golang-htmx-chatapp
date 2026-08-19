package repository

import (
	"errors"

	"gorm.io/gorm"
)

// ErrNotFound is the repository layer's own not-found signal. Callers
// (services) check against this, never against gorm's error types
// directly, so gorm stays an implementation detail of this package.
var ErrNotFound = errors.New("record not found")

type Repositories struct {
	UserRepo      *UserRepository
	GroupRepo     *GroupRepository
	UserGroupRepo *UserGroupRepository
	MessageRepo   *MessageRepository
}

// Deps holds the dependencies needed to construct the repository layer.
type Deps struct {
	DB *gorm.DB
}

func NewRepositories(deps Deps) *Repositories {
	return &Repositories{
		UserRepo:      NewUserRepository(deps.DB),
		GroupRepo:     NewGroupRepository(deps.DB),
		UserGroupRepo: NewUserGroupRepository(deps.DB),
		MessageRepo:   NewMessageRepository(deps.DB),
	}
}
