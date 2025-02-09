package repository

import (
	"context"
	db "templates_new/internal/client"
	"templates_new/internal/models"
	"templates_new/pkg/token"

	"github.com/rs/zerolog"
)

type Autorization interface {
	SignIn(ctx context.Context, user *models.User) (*models.User, error)
	LogIn(ctx context.Context, user *models.User) (*models.User, error)
	CreateSession(ctx context.Context, user *models.User, userClaims *token.UserClaims, refreshToken string) (string, error)
	GetSessionById(ctx context.Context, id string) (*models.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

type Repository struct {
	Autorization
}

func NewRepository(db db.Client, log zerolog.Logger) *Repository {
	return &Repository{Autorization: NewAuthRepository(db, log)}
}
