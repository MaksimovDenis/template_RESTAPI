package server

import (
	"context"
	"fmt"
	db "templates_new/internal/client"
	"templates_new/internal/repository"
	"templates_new/pkg/protocol/oapi"

	"github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ServiceRepository {
	return &repo{db: db}
}

func (rep *repo) SignIn(ctx context.Context, user *oapi.CreateUserJSONBody) (*oapi.User, error) {
	builder := squirrel.Insert("user").
		PlaceholderFormat(squirrel.Dollar).
		Columns("username", "email", "password", "is_admin").
		Values(user.Username, user.Email, user.Password, user.IsAdmin)

	query, args, err := builder.ToSql()
	fmt.Println(query, args)
	if err != nil {
		return nil, err
	}

	queryStruct := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRow: query,
	}

	res := &oapi.User{}

	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Username, &res.Email, res.IsAdmin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return res, nil
}

func (rep *repo) LogIn(ctx context.Context, user *oapi.User) (*oapi.CreateUserJSONBody, error) {
	builder := squirrel.Select("username", "email", "password", "is_admin").
		PlaceholderFormat(squirrel.Dollar).
		From("users").
		Where(squirrel.Eq{"username": user.Username})

	query, args, err := builder.ToSql()
	fmt.Println(query, args)
	if err != nil {
		return nil, err
	}

	queryStruct := db.Query{
		Name:     "user_repository.Login",
		QueryRow: query,
	}

	res := &oapi.CreateUserJSONBody{}

	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Username, &res.Email, &res.Password, &res.IsAdmin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return res, nil
}
