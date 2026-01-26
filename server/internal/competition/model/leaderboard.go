package model

type LeaderboardEntry struct {
	TradingAccountLogin int64
	Rank                int32
	Username            string
	AccountSize         float64
	Profit              float64
	Equity              float64
	GainPercent         float64
}
