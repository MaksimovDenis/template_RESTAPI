package handler

import (
	"fmt"
	"net/http"
	"strings"
	"templates_new/pkg/token"

	"github.com/gin-gonic/gin"
)

func GetAuthMiddlewareFunc(tokenMaker *token.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// read the authorization header
		// verify the token
		claims, err := verifyClaimsFromAuthHeader(ctx, *tokenMaker)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("user", claims)
		ctx.Next()
	}
}

func GetAdminMiddlewareFunc(tokenMaker *token.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// read the authorization header
		// verify the token
		claims, err := verifyClaimsFromAuthHeader(ctx, *tokenMaker)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if !claims.IsAdmin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user is not an admin"})
			return
		}

		ctx.Set("user", claims)
		ctx.Next()
	}
}

func verifyClaimsFromAuthHeader(ctx *gin.Context, tokenMaker token.JWTMaker) (*token.UserClaims, error) {
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("autorizatoin header is missing")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("invalid autorization header")
	}

	token := fields[1]
	claims, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
