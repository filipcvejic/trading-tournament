package model

import "github.com/google/uuid"

type CompetitionMember struct {
	CompetitionID uuid.UUID
	TradingLogin  int64
	AccountSize   float64
}
