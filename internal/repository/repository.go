package repository

import (
	"context"
	"templates_new/pkg/protocol/oapi"
)

type ServiceRepository interface {
	SignIn(ctx context.Context, user *oapi.CreateUserJSONBody) (*oapi.User, error)
	LogIn(ctx context.Context, user *oapi.User) (*oapi.CreateUserJSONBody, error)
}
