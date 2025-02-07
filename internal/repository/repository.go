package repository

import (
	"context"
	db "templates_new/internal/client"
	"templates_new/pkg/protocol/oapi"
)

type Autorization interface {
	SignIn(ctx context.Context, user *oapi.SignInJSONBody) (*oapi.User, error)
	LogIn(ctx context.Context, user *oapi.User) (*oapi.SignInJSONBody, error)
}

type Repository struct {
	Autorization
}

func NewRepository(db db.Client) *Repository {
	return &Repository{Autorization: NewAuthRepository(db)}
}
