package repository

import "gorm.io/gorm"

type Repositories struct {
	UserRepo *UserRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo: NewUserRepository(db),
	}
}
