package service

import (
	db "templates_new/internal/client"
	"templates_new/internal/repository"
	"templates_new/pkg/token"

	"github.com/rs/zerolog"
)

type Service struct {
	Authorization Authorization
	Server        Server
}

func NewService(repos repository.Repository, txManager db.TxManager, token token.JWTMaker, log zerolog.Logger) *Service {
	return &Service{
		Authorization: newAuthService(repos, txManager, token, log),
		Server:        newServerService(),
	}
}
