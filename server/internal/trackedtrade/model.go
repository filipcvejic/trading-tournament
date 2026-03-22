package trackedtrade

import "time"

type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

type TrackedTrade struct {
	PositionID int64
	Symbol     string
	Side       Side
	OpenPrice  float64
	StopLoss   *float64
	OpenedAt   time.Time
	ClosedAt   *time.Time
}
