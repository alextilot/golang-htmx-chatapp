package repository

import "gorm.io/gorm"

type Repositories struct {
	UserRepo  *UserRepository
	GroupRepo *GroupRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:  NewUserRepository(db),
		GroupRepo: NewGroupRepository(db),
	}
}
