package user

import "errors"

var (
	ErrUserNotFound                 = errors.New("user not found")
	ErrInvalidEmail                 = errors.New("invalid email")
	ErrInvalidUsername              = errors.New("invalid username")
	ErrInvalidDiscordUsername       = errors.New("invalid discord username")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrEmailAlreadyExists           = errors.New("email already exists")
	ErrUsernameAlreadyExists        = errors.New("username already exists")
	ErrDiscordUsernameAlreadyExists = errors.New("discord username already exists")
)
