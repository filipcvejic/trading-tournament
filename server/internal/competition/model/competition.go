package model

import (
	"github.com/google/uuid"
	"time"
)

type Competition struct {
	ID        uuid.UUID
	Name      string
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedAt time.Time
}
