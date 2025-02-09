package models

import (
	"time"
)

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"Email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserRes struct {
	User                  User      `json:"user"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	SessionId             int       `json:"session_id"`
}

type Session struct {
	Id           int       `json:"id"`
	UserEmail    string    `json:"user_email"`
	RefreshToken string    `json:"refresh_token"`
	IsRevoked    bool      `json:"is_revoked"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RenewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
