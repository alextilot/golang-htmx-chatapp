package service

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
)

type Services struct {
	UserService *UserService
	// ChatGroupService   *ChatGroupService
	// ChatMessageService *ChatMessageService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		UserService: NewUserService(repos.UserRepo),
		// ChatGroupService:   NewChatGroupService(repos.ChatGroupRepo),
		// ChatMessageService: NewChatMessageService(repos.ChatMessageRepo),
	}
}
