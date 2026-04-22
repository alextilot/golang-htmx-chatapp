package service

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
)

type Services struct {
	UserService  *UserService
	GroupService *GroupService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		UserService:  NewUserService(repos.UserRepo),
		GroupService: NewGroupService(repos.GroupRepo, repos.UserGroupRepo),
	}
}
