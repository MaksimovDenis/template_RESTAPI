package repository

import (
	db "templates_new/internal/client"

	"github.com/rs/zerolog"
)

type Repository struct {
	Authorization Authorization
}

func NewRepository(db db.Client, log zerolog.Logger) *Repository {
	return &Repository{Authorization: newAuthRepository(db, log)}
}
