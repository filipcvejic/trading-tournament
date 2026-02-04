package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,min=3,max=20,no_whitespace"`
	DiscordUsername string `json:"discordUsername" validate:"required,min=2,max=32,no_whitespace"`
	Password        string `json:"password" validate:"required,password_strong"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type AuthResponse struct {
	AccessToken string `json:"AccessToken"`
}

type AuthWithRefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
