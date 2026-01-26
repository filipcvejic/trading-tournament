package mapper

import (
	"github.com/filipcvejic/trading_tournament/db/sqlc"
	"github.com/filipcvejic/trading_tournament/internal/competition/dto"
	"github.com/filipcvejic/trading_tournament/internal/competition/model"
)

func LeaderboardFromDB(rows []sqlc.GetCompetitionLeaderboardRow) []model.LeaderboardEntry {
	entries := make([]model.LeaderboardEntry, 0, len(rows))

	for _, r := range rows {
		entries = append(entries, model.LeaderboardEntry{
			TradingAccountLogin: r.TradingAccountLogin,
			Rank:                r.Rank,
			Username:            r.Username,
			AccountSize:         r.AccountSize,
			Profit:              r.Profit,
			Equity:              r.Equity,
			GainPercent:         r.GainPercent,
		})
	}

	return entries
}

func LeaderboardToDTO(entries []model.LeaderboardEntry) []dto.LeaderboardEntryResponse {
	out := make([]dto.LeaderboardEntryResponse, 0, len(entries))

	for _, e := range entries {
		out = append(out, dto.LeaderboardEntryResponse{
			TradingAccountLogin: e.TradingAccountLogin,
			Rank:                e.Rank,
			Username:            e.Username,
			AccountSize:         e.AccountSize,
			Profit:              e.Profit,
			Equity:              e.Equity,
			GainPercent:         e.GainPercent,
		})
	}

	return out
}
