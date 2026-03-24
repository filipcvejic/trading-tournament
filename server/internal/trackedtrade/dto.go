package trackedtrade

import "time"

type EventType string

const (
	EventTypeOpen  EventType = "OPEN"
	EventTypeClose EventType = "CLOSE"
)

type IngestTrackedTradeEventRequest struct {
	EventType  EventType  `json:"eventType" validate:"required,oneof=OPEN CLOSE"`
	PositionID int64      `json:"positionId" validate:"required"`
	Symbol     string     `json:"symbol,omitempty"`
	Side       string     `json:"side,omitempty"`
	OpenPrice  float64    `json:"openPrice,omitempty"`
	StopLoss   *float64   `json:"stopLoss,omitempty"`
	Volume     float64    `json:"volume,omitempty"`
	OpenedAt   *time.Time `json:"openedAt,omitempty"`
	ClosedAt   *time.Time `json:"closedAt,omitempty"`
}
