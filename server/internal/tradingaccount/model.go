package tradingaccount

import (
	"github.com/google/uuid"
	"time"
)

type TradingAccount struct {
	Login     int64
	UserID    uuid.UUID
	Broker    string
	CreatedAt time.Time
}
