package dto

import (
	"time"
)

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

type InsertTradesRequest struct {
	TradingAccountLogin int64      `json:"accountId"`
	Trades              []TradeDTO `json:"trades"`
}
