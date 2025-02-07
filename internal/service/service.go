package service

import (
	"context"
	db "templates_new/internal/client"
	"templates_new/internal/repository"
	"templates_new/pkg/protocol/oapi"

	"github.com/gin-gonic/gin"
)

type Server interface {
	CheckService(ctx *gin.Context) error
}

type Autorization interface {
	SignIn(ctx context.Context, user *oapi.SignInJSONBody) (*oapi.User, error)
	LogIn(ctx context.Context, user *oapi.User) (*oapi.SignInJSONBody, error)
}

type Service struct {
	Autorization
	Server
}

func NewService(repos repository.Repository, txManager db.TxManager) *Service {
	return &Service{
		Autorization: NewAuthService(repos, txManager),
		Server:       NewServerService(),
	}
}
