package tradingaccount

import (
	"context"
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	login int64,
	userID uuid.UUID,
	broker string,
	investorPasswordEncrypted string,
) (TradingAccount, error) {
	if login <= 0 {
		return TradingAccount{}, ErrInvalidLogin
	}
	if userID == uuid.Nil {
		return TradingAccount{}, user.ErrNotFound
	}
	if broker == "" || strings.TrimSpace(broker) != broker {
		return TradingAccount{}, ErrInvalidBroker
	}
	if investorPasswordEncrypted == "" {
		return TradingAccount{}, ErrInvalidInvestorPassword
	}

	acc, err := s.repo.Create(ctx, login, userID, broker, investorPasswordEncrypted)
	if err == nil {
		return acc, nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return TradingAccount{}, ErrLoginTaken
		case "23503":
			return TradingAccount{}, user.ErrNotFound
		}
	}

	return TradingAccount{}, err
}

func (s *Service) GetByLogin(ctx context.Context, login int64) (TradingAccount, error) {
	if login <= 0 {
		return TradingAccount{}, ErrInvalidLogin
	}
	return s.repo.GetByLogin(ctx, login)
}

func (s *Service) GetTradeHistory(ctx context.Context, login int64) (TradeHistoryResponse, error) {
	if login <= 0 {
		return TradeHistoryResponse{}, ErrInvalidLogin
	}

	username, trades, err := s.repo.GetTradeHistory(ctx, login)
	if err != nil {
		return TradeHistoryResponse{}, err
	}

	return TradeHistoryResponse{
		Username: username,
		Trades:   trades,
	}, nil
}
