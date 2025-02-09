package service

import (
	"context"
	"errors"
	db "templates_new/internal/client"
	"templates_new/internal/models"
	"templates_new/internal/repository"
	"templates_new/pkg/token"
	"templates_new/pkg/util"
	"time"

	"github.com/rs/zerolog"
)

const (
	durationRefreshToken time.Duration = 24 * time.Hour

	// TODO: JUST FOR TESTING (IN ORDINARY CASE 15 MIN)
	durationAccessToken time.Duration = 3 * time.Hour
)

type AuthService struct {
	appRepository repository.Repository
	txManager     db.TxManager
	token         token.JWTMaker
	log           zerolog.Logger
}

func NewAuthService(
	appRepository repository.Repository,
	txManager db.TxManager,
	token token.JWTMaker,
	log zerolog.Logger,
) *AuthService {
	return &AuthService{
		appRepository: appRepository,
		txManager:     txManager,
		token:         token,
		log:           log,
	}
}

func (auth *AuthService) SignIn(ctx context.Context, user *models.User) (*models.User, error) {
	hashed, err := util.HashPassword(user.Password)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to hash password")
		return nil, err
	}

	user.Password = hashed

	userCreate, err := auth.appRepository.SignIn(ctx, user)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to create user")
		return nil, err
	}

	return userCreate, nil
}

func (auth *AuthService) LogIn(ctx context.Context, user *models.User) (*models.UserRes, error) {
	userStorage, err := auth.appRepository.Autorization.LogIn(ctx, user)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to retrieve user from storage")
		return nil, err
	}

	if err = util.CheckPassword(user.Password, userStorage.Password); err != nil {
		auth.log.Error().Err(err).Msg("password mismatch")
		return nil, err
	}

	accessToken, accessClaims, err := auth.token.CreateToken(int64(userStorage.Id), userStorage.Email, userStorage.IsAdmin, durationAccessToken)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to create access token")
		return nil, err
	}

	refreshToken, refreshClaims, err := auth.token.CreateToken(int64(userStorage.Id), userStorage.Email, userStorage.IsAdmin, durationRefreshToken)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to create refresh token")
		return nil, err
	}

	sessionId, err := auth.appRepository.CreateSession(ctx, userStorage, refreshClaims, refreshToken)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to create session")
		return nil, err
	}

	res := &models.UserRes{
		User: models.User{
			UserName: userStorage.UserName,
			Email:    userStorage.Email,
			IsAdmin:  userStorage.IsAdmin,
		},
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		SessionId:             sessionId,
	}

	return res, nil
}

func (auth *AuthService) LogOut(ctx context.Context, id string) error {
	if err := auth.appRepository.DeleteSession(ctx, id); err != nil {
		auth.log.Error().Err(err).Msg("failed to delete session")
		return err
	}
	return nil
}

func (auth *AuthService) RenewAccessToken(ctx context.Context, refreshToken string) (*token.UserClaims, string, error) {
	refreshClaims, err := auth.token.VerifyToken(refreshToken)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to verify refresh token")
		return nil, "", err
	}

	session, err := auth.appRepository.GetSessionById(ctx, refreshClaims.RegisteredClaims.ID)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to retrieve session")
		return nil, "", err
	}

	if session.IsRevoked {
		auth.log.Error().Msg("session is revoked")
		return nil, "", errors.New("session revoked")
	}

	if session.UserEmail != refreshClaims.Email {
		auth.log.Error().Msg("session email mismatch")
		return nil, "", errors.New("invalid session")
	}

	accessToken, accessClaims, err := auth.token.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.IsAdmin, durationAccessToken)
	if err != nil {
		auth.log.Error().Err(err).Msg("failed to create new access token")
		return nil, "", err
	}

	return accessClaims, accessToken, nil
}
