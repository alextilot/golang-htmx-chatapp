package repository

import "gorm.io/gorm"

type Repositories struct {
	UserRepo        *UserRepository
	ChatGroupRepo   *ChatGroupRepository
	ChatMessageRepo *ChatMessageRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:        NewUserRepository(db),
		ChatGroupRepo:   NewChatGroupRepository(db),
		ChatMessageRepo: NewChatMessageRepository(db),
	}
}
