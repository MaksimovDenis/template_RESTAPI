package service

import (
	"context"
	db "templates_new/internal/client"
	"templates_new/internal/repository"
	"templates_new/pkg/protocol/oapi"
	"templates_new/pkg/util"
)

type AuthService struct {
	appRepository repository.Repository
	txManager     db.TxManager
}

func NewAuthService(
	appRepository repository.Repository,
	txManager db.TxManager,
) *AuthService {
	return &AuthService{
		appRepository: appRepository,
		txManager:     txManager,
	}
}

func (auth *AuthService) SignIn(ctx context.Context, user *oapi.SignInJSONBody) (*oapi.User, error) {
	hashed, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashed

	userCreate, err := auth.appRepository.SignIn(ctx, user)
	if err != nil {
		return nil, err
	}

	return userCreate, nil
}

func (auth *AuthService) LogIn(ctx context.Context, user *oapi.User) (*oapi.SignInJSONBody, error) {
	return nil, nil
}
