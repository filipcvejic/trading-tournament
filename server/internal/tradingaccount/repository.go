package tradingaccount

import (
	"context"
	"database/sql"
	"errors"
	"github.com/filipcvejic/trading_tournament/db"
	"github.com/filipcvejic/trading_tournament/db/sqlc"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, login int64, userID uuid.UUID, broker, investorPasswordEncrypted string) (TradingAccount, error)
	GetByLogin(ctx context.Context, login int64) (TradingAccount, error)
	GetTradeHistory(ctx context.Context, login int64) (username string, trades []TradeDTO, err error)
}

type PostgresRepository struct {
	db *db.DB
}

func NewPostgresRepository(database *db.DB) *PostgresRepository {
	return &PostgresRepository{db: database}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	login int64,
	userID uuid.UUID,
	broker string,
	investorPasswordEncrypted string,
) (TradingAccount, error) {
	row, err := r.db.Query.CreateTradingAccount(ctx, sqlc.CreateTradingAccountParams{
		Login:                     login,
		UserID:                    userID,
		Broker:                    broker,
		InvestorPasswordEncrypted: investorPasswordEncrypted,
	})
	if err != nil {
		return TradingAccount{}, err
	}

	return TradingAccount{
		Login:     row.Login,
		UserID:    row.UserID,
		Broker:    row.Broker,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *PostgresRepository) GetByLogin(ctx context.Context, login int64) (TradingAccount, error) {
	row, err := r.db.Query.GetTradingAccountByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TradingAccount{}, ErrNotFound
		}
		return TradingAccount{}, err
	}

	return TradingAccount{
		Login:     row.Login,
		UserID:    row.UserID,
		Broker:    row.Broker,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *PostgresRepository) GetTradeHistory(
	ctx context.Context,
	login int64,
) (string, []TradeDTO, error) {
	username, err := r.db.Query.GetUsernameByTradingAccountLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, ErrNotFound
		}
		return "", nil, err
	}

	rows, err := r.db.Query.ListTradesByAccountLogin(ctx, login)
	if err != nil {
		return "", nil, err
	}

	trades := make([]TradeDTO, 0, len(rows))
	for _, row := range rows {
		trades = append(trades, TradeDTO{
			PositionID: row.PositionID,
			Symbol:     row.Symbol,
			Side:       row.Side,
			Volume:     row.Volume,
			OpenTime:   row.OpenTime,
			CloseTime:  row.CloseTime,
			OpenPrice:  row.OpenPrice,
			ClosePrice: row.ClosePrice,
			Profit:     row.Profit,
			Commission: row.Commission,
			Swap:       row.Swap,
		})
	}

	return username, trades, nil
}
