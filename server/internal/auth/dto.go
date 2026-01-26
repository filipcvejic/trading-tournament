package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	DiscordUsername string `json:"discordUsername"`
	Password        string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
