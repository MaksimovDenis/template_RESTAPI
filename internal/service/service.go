package service

import (
	"context"
	db "templates_new/internal/client"
	"templates_new/internal/models"
	"templates_new/internal/repository"
	"templates_new/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server interface {
	CheckService(ctx *gin.Context) error
}

type Autorization interface {
	SignIn(ctx context.Context, user *models.User) (*models.User, error)
	LogIn(ctx context.Context, user *models.User) (*models.UserRes, error)
	LogOut(ctx context.Context, id int) error
	RenewAccessToken(ctx context.Context, refreshToken string) (*token.UserClaims, string, error)
}

type Service struct {
	Autorization
	Server
}

func NewService(repos repository.Repository, txManager db.TxManager, token token.JWTMaker, log zerolog.Logger) *Service {
	return &Service{
		Autorization: NewAuthService(repos, txManager, token, log),
		Server:       NewServerService(),
	}
}
