package repository

import "gorm.io/gorm"

type Repositories struct {
	UserRepo      *UserRepository
	GroupRepo     *GroupRepository
	UserGroupRepo *UserGroupRepository
	MessageRepo   *MessageRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:      NewUserRepository(db),
		GroupRepo:     NewGroupRepository(db),
		UserGroupRepo: NewUserGroupRepository(db),
		MessageRepo:   NewMessageRepository(db),
	}
}
