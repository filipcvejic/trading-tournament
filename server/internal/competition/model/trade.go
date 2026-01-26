package model

import (
	"github.com/google/uuid"
	"time"
)

type Trade struct {
	TradingAccountLogin int64
	CompetitionID       uuid.UUID
	PositionID          int64
	Symbol              string
	Side                string
	Volume              float64
	OpenTime            time.Time
	CloseTime           time.Time
	OpenPrice           float64
	ClosePrice          float64
	Profit              float64
	Commission          float64
	Swap                float64
}
