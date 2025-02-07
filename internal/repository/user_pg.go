package repository

import (
	"context"
	"fmt"
	db "templates_new/internal/client"
	"templates_new/pkg/protocol/oapi"

	"github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthRepo struct {
	db db.Client
}

func NewAuthRepository(db db.Client) *AuthRepo {
	return &AuthRepo{db: db}
}

func (rep *AuthRepo) SignIn(ctx context.Context, user *oapi.SignInJSONBody) (*oapi.User, error) {
	builder := squirrel.Insert("users").
		PlaceholderFormat(squirrel.Dollar).
		Columns("username", "email", "password", "is_admin").
		Values(user.Username, user.Email, user.Password, user.IsAdmin).
		Suffix("RETURNING username, email, is_admin")

	query, args, err := builder.ToSql()
	fmt.Println(query, args)
	if err != nil {
		return nil, err
	}

	queryStruct := db.Query{
		Name:     "user_repository.SignIn",
		QueryRow: query,
	}

	res := &oapi.User{}

	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args...).Scan(&res.Username, &res.Email, &res.IsAdmin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return res, nil
}

func (rep *AuthRepo) LogIn(ctx context.Context, user *oapi.User) (*oapi.SignInJSONBody, error) {
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
		Name:     "user_repository.LogIn",
		QueryRow: query,
	}

	res := &oapi.SignInJSONBody{}

	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args).
		Scan(&res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return res, nil
}
