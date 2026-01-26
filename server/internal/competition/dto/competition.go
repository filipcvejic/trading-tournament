package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateCompetitionRequest struct {
	Name     string    `json:"name"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}

type CompetitionResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}

type JoinCompetitionRequest struct {
	Login            int64  `json:"login"`
	Broker           string `json:"broker"`
	InvestorPassword string `json:"investorPassword"`
}

type UpdateAccountSizeRequest struct {
	TradingAccountLogin int64   `json:"tradingAccountLogin"`
	AccountSize         float64 `json:"accountSize"`
}

type CompetitionUserStateResponse struct {
	HasRequestedAccount bool `json:"hasRequestedAccount"`
	HasJoined           bool `json:"hasJoined"`
}
