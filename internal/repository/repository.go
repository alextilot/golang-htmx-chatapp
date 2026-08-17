package repository

import "gorm.io/gorm"

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
