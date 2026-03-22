package trackedtrade

import (
	"context"
	"github.com/filipcvejic/trading_tournament/db"
	"time"

	"github.com/filipcvejic/trading_tournament/db/sqlc"
)

type Repository interface {
	Create(ctx context.Context, trade TrackedTrade) error
	Close(ctx context.Context, positionID int64, closedAt time.Time) error
	List(ctx context.Context) ([]TrackedTrade, error)
}

type PostgresRepository struct {
	db *db.DB
}

func NewPostgresRepository(database *db.DB) *PostgresRepository {
	return &PostgresRepository{db: database}
}

func (r *PostgresRepository) Create(ctx context.Context, trade TrackedTrade) error {
	return r.db.Query.CreateTrackedTrade(ctx, sqlc.CreateTrackedTradeParams{
		PositionID: trade.PositionID,
		Symbol:     trade.Symbol,
		Side:       string(trade.Side),
		OpenPrice:  trade.OpenPrice,
		StopLoss:   trade.StopLoss,
		OpenedAt:   trade.OpenedAt,
	})
}

func (r *PostgresRepository) Close(ctx context.Context, positionID int64, closedAt time.Time) error {
	return r.db.Query.CloseTrackedTrade(ctx, sqlc.CloseTrackedTradeParams{
		PositionID: positionID,
		ClosedAt:   &closedAt,
	})
}

func (r *PostgresRepository) List(ctx context.Context) ([]TrackedTrade, error) {
	rows, err := r.db.Query.ListTrackedTrades(ctx)
	if err != nil {
		return nil, err
	}

	trades := make([]TrackedTrade, 0, len(rows))
	for _, row := range rows {
		trades = append(trades, TrackedTrade{
			PositionID: row.PositionID,
			Symbol:     row.Symbol,
			Side:       Side(row.Side),
			OpenPrice:  row.OpenPrice,
			StopLoss:   row.StopLoss,
			OpenedAt:   row.OpenedAt,
			ClosedAt:   row.ClosedAt,
		})
	}

	return trades, nil
}
