package dto

type LeaderboardEntryResponse struct {
	TradingAccountLogin int64   `json:"tradingAccountLogin"`
	Rank                int32   `json:"rank"`
	Username            string  `json:"username"`
	AccountSize         float64 `json:"accountSize"`
	Profit              float64 `json:"profit"`
	Equity              float64 `json:"equity"`
	GainPercent         float64 `json:"gainPercent"`
}
