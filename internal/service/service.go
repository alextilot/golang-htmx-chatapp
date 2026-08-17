package service

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
)

type Services struct {
	UserService  *UserService
	GroupService *GroupService
}

// Deps holds the dependencies needed to construct the service layer.
type Deps struct {
	Repos *repository.Repositories
}

func NewServices(deps Deps) *Services {
	return &Services{
		UserService:  NewUserService(deps.Repos.UserRepo),
		GroupService: NewGroupService(deps.Repos.GroupRepo, deps.Repos.UserGroupRepo, deps.Repos.MessageRepo),
	}
}
