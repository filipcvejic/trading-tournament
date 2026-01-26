package mapper

import (
	"github.com/filipcvejic/trading_tournament/db/sqlc"
	"github.com/filipcvejic/trading_tournament/internal/competition/model"
)

func CompetitionFromDB(row sqlc.Competition) model.Competition {
	return model.Competition{
		ID:        row.ID,
		Name:      row.Name,
		StartsAt:  row.StartsAt,
		EndsAt:    row.EndsAt,
		CreatedAt: row.CreatedAt,
	}
}

//func CompetitionToDTO(c model.Competition) dto.CompetitionResponse {
//	return dto.CompetitionResponse{
//		ID:       c.ID.String(),
//		Name:     c.Name,
//		StartsAt: c.StartsAt.Format(time.RFC3339),
//		EndsAt:   c.EndsAt.Format(time.RFC3339),
//	}
//}
