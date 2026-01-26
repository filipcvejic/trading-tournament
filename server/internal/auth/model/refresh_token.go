package model

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}
