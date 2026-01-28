package tradingaccount

import (
	"github.com/google/uuid"
	"time"
)

type CreateTradingAccountRequest struct {
	Login            int64     `json:"login" validate:"required,min=10000"`
	UserID           uuid.UUID `json:"userId"`
	Broker           string    `json:"broker" validate:"required,min=2"`
	InvestorPassword string    `json:"investorPassword" validate:"required,min=5"`
}

type TradingAccountResponse struct {
	Login     int64     `json:"login"`
	UserID    uuid.UUID `json:"userId"`
	Broker    string    `json:"broker"`
	CreatedAt time.Time `json:"createdAt"`
}

type TradeHistoryResponse struct {
	Username string     `json:"username"`
	Trades   []TradeDTO `json:"trades"`
}

type TradeDTO struct {
	PositionID int64     `json:"positionId"`
	Symbol     string    `json:"symbol"`
	Side       string    `json:"side"`
	Volume     float64   `json:"volume"`
	OpenTime   time.Time `json:"openTime"`
	CloseTime  time.Time `json:"closeTime"`
	OpenPrice  float64   `json:"openPrice"`
	ClosePrice float64   `json:"closePrice"`
	Profit     float64   `json:"profit"`
	Commission float64   `json:"commission"`
	Swap       float64   `json:"swap"`
}
