package user

import "time"

type CreateUserRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	DiscordUsername string `json:"discordUsername"`
	Password        string `json:"password"`
}

type UserResponse struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	DiscordUsername string    `json:"discordUsername"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
