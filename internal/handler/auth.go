package handler

import (
	"net/http"
	"templates_new/internal/models"
	"templates_new/pkg/protocol/oapi"
	"templates_new/pkg/token"

	"github.com/gin-gonic/gin"
)

func (hdl *Handler) SignIn(ctx *gin.Context) {
	var userReq oapi.UserReq

	if err := ctx.BindJSON(&userReq); err != nil {
		hdl.log.Error().Err(err).Msg("failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse body request"})

		return
	}

	user := &models.User{
		UserName: userReq.Username,
		Email:    userReq.Email,
		Password: userReq.Password,
		IsAdmin:  userReq.IsAdmin,
	}

	service, err := hdl.appService.Authorization.SignIn(ctx, user)
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to sign in user")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to sign in user"})

		return
	}

	userRes := &oapi.UserRes{
		Username: service.UserName,
		Email:    service.Email,
		IsAdmin:  service.IsAdmin,
	}

	ctx.JSON(http.StatusCreated, userRes)
}

func (hdl *Handler) LogIn(ctx *gin.Context) {
	var loginReq oapi.LoginUserReq

	if err := ctx.BindJSON(&loginReq); err != nil {
		hdl.log.Error().Err(err).Msg("failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse body request"})

		return
	}

	user := &models.User{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	}

	res, err := hdl.appService.Authorization.LogIn(ctx, user)
	if err != nil {
		hdl.log.Error().Err(err).Msg("invalid email or password")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})

		return
	}

	loginRes := &oapi.LoginUserRes{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		SessionId:    res.SessionId,
		User: oapi.UserRes{
			Username: res.User.UserName,
			Email:    res.User.Email,
			IsAdmin:  res.User.IsAdmin,
		},
	}

	ctx.JSON(http.StatusOK, loginRes)
}

func (hdl *Handler) LogOut(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	if !ok {
		hdl.log.Error().Msg("user claims not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

		return
	}

	sessionId := claims.(*token.UserClaims).RegisteredClaims.ID

	if err := hdl.appService.Authorization.LogOut(ctx, sessionId); err != nil {
		hdl.log.Error().Err(err).Msg("failed to log out user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (hdl *Handler) RenewAccessToken(ctx *gin.Context) {
	var refreshToken oapi.RenewAccessTokenReq

	if err := ctx.BindJSON(&refreshToken); err != nil {
		hdl.log.Error().Err(err).Msg("invalid refresh token request format")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})

		return
	}

	userClaims, accessToken, err := hdl.appService.Authorization.RenewAccessToken(ctx, refreshToken.RefreshToken)
	if err != nil {
		hdl.log.Error().Err(err).Msg("failed to renew access token")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get access token"})

		return
	}

	res := &oapi.RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: oapi.AccessTokenExpiresAt(*userClaims.ExpiresAt),
	}

	ctx.JSON(http.StatusOK, res)
}
