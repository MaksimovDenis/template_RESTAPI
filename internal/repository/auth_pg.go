package repository

import (
	"context"
	db "templates_new/internal/client"
	"templates_new/internal/models"
	"templates_new/pkg/token"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthRepo struct {
	db  db.Client
	log zerolog.Logger
}

func NewAuthRepository(db db.Client, log zerolog.Logger) *AuthRepo {
	return &AuthRepo{
		db:  db,
		log: log,
	}
}

func (rep *AuthRepo) SignIn(ctx context.Context, user *models.User) (*models.User, error) {
	builder := squirrel.Insert("users").
		PlaceholderFormat(squirrel.Dollar).
		Columns("username", "email", "password", "is_admin").
		Values(user.UserName, user.Email, user.Password, user.IsAdmin).
		Suffix("RETURNING username, email, is_admin")

	query, args, err := builder.ToSql()
	if err != nil {
		rep.log.Error().Err(err).Msg("SignIn: failed to build SQL query")
		return nil, err
	}

	queryStruct := db.Query{
		Name:     "user_repository.SignIn",
		QueryRow: query,
	}

	res := &models.User{}

	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args...).Scan(&res.UserName, &res.Email, &res.IsAdmin)
	if err != nil {
		rep.log.Error().Err(err).Msg("SignIn: failed to execute query")
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return res, nil
}

func (rep *AuthRepo) LogIn(ctx context.Context, user *models.User) (*models.User, error) {
	builder := squirrel.Select("username", "email", "password", "is_admin").
		PlaceholderFormat(squirrel.Dollar).
		From("users").
		Where(squirrel.Eq{"email": user.Email})

	query, args, err := builder.ToSql()
	if err != nil {
		rep.log.Error().Err(err).Msg("LogIn: failed to build SQL query")
		return nil, err
	}

	queryStruct := db.Query{
		Name:     "user_repository.LogIn",
		QueryRow: query,
	}

	res := &models.User{}

	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.UserName, &res.Email, &res.Password, &res.IsAdmin)
	if err != nil {
		rep.log.Error().Err(err).Msg("LogIn: failed to execute query")
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return res, nil
}

func (rep *AuthRepo) CreateSession(ctx context.Context,
	user *models.User,
	userClaims *token.UserClaims,
	refreshToken string) (string, error) {
	builder := squirrel.Insert("sessions").
		PlaceholderFormat(squirrel.Dollar).
		Columns("id", "user_email", "refresh_token", "expires_at").
		Values(userClaims.RegisteredClaims.ID, user.Email, refreshToken, userClaims.ExpiresAt.Time).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		rep.log.Error().Err(err).Msg("CreateSession: failed to build SQL query")
		return "", err
	}

	queryStruct := db.Query{
		Name:     "user_repository.CreateSession",
		QueryRow: query,
	}

	var sessionId string

	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args...).Scan(&sessionId)
	if err != nil {
		rep.log.Error().Err(err).Msg("CreateSession: failed to execute query")
		return "", status.Errorf(codes.Internal, "Internal server error")
	}

	return sessionId, nil
}

func (rep *AuthRepo) DeleteSession(ctx context.Context, id string) error {
	builder := squirrel.Delete("sessions").
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		rep.log.Error().Err(err).Msg("DeleteSession: failed to build SQL query")
		return err
	}

	queryStruct := db.Query{
		Name:     "user_repository.DeleteSession",
		QueryRow: query,
	}

	res, err := rep.db.DB().ExecContext(ctx, queryStruct, args...)
	if err != nil {
		rep.log.Error().Err(err).Msg("DeleteSession: failed to execute query")
		return status.Errorf(codes.Internal, "Internal server error")
	}

	rowAffected := res.RowsAffected()
	if rowAffected == 0 {
		rep.log.Error().Msg("DeleteSession: no rows affected")
		return status.Errorf(codes.NotFound, "Session not found")
	}

	return nil
}

func (rep *AuthRepo) GetSessionById(ctx context.Context, id string) (*models.Session, error) {
	builder := squirrel.Select("*").
		PlaceholderFormat(squirrel.Dollar).
		From("sessions").
		Where(squirrel.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		rep.log.Error().Err(err).Msg("GetSessionById: failed to build SQL query")
		return nil, err
	}

	queryStruct := db.Query{
		Name:     "user_repository.GetSessionById",
		QueryRow: query,
	}

	res := &models.Session{}

	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args...).
		Scan(&res.Id, &res.UserEmail, &res.RefreshToken, &res.IsRevoked, &res.CreatedAt, &res.ExpiresAt)
	if err != nil {
		rep.log.Error().Err(err).Msg("GetSessionById: failed to execute query")
		return nil, err
	}

	return res, nil
}
