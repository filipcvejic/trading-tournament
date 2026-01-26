package auth

import "errors"

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmailInUse           = errors.New("email already in use")
	ErrUsernameInUse        = errors.New("username already in use")
	ErrDiscordUsernameInUse = errors.New("discord username already in use")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("token expired")
	ErrUnauthorized         = errors.New("unauthorized")
)
