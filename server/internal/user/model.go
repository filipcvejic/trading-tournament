package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Email           string
	Username        string
	DiscordUsername string
	PasswordHash    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
