package competition

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/filipcvejic/trading_tournament/db"
	"github.com/filipcvejic/trading_tournament/db/sqlc"
	"github.com/filipcvejic/trading_tournament/internal/competition/mapper"
	"github.com/filipcvejic/trading_tournament/internal/competition/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, c model.Competition) error
	GetByID(ctx context.Context, id uuid.UUID) (model.Competition, error)
	JoinWithTradingAccount(ctx context.Context, competitionID uuid.UUID, userID uuid.UUID, login int64, broker string, investorPasswordEncrypted string) error
	UpdateAccountSize(ctx context.Context, competitionID uuid.UUID, login int64, accountSize float64) error
	GetMemberAccountSize(ctx context.Context, competitionID uuid.UUID, login int64) (float64, error)
	GetLeaderboard(ctx context.Context, competitionID uuid.UUID, limit, offset int32) ([]model.LeaderboardEntry, error)
	InsertTrades(ctx context.Context, competitionID uuid.UUID, login int64, trades []model.Trade) error
	GetUserCompetitionState(ctx context.Context, userID, competitionID uuid.UUID) (sqlc.GetCompetitionUserStateRow, error)
	GetCurrent(ctx context.Context) (sqlc.Competition, error)
	CreateAccountRequest(ctx context.Context, userID, competitionID uuid.UUID) error
}

type PostgresRepository struct {
	db *db.DB
}

func NewPostgresRepository(database *db.DB) *PostgresRepository {
	return &PostgresRepository{db: database}
}

func (r *PostgresRepository) Create(ctx context.Context, c model.Competition) error {
	_, err := r.db.Query.CreateCompetition(ctx, sqlc.CreateCompetitionParams{
		ID:       c.ID,
		Name:     c.Name,
		StartsAt: c.StartsAt,
		EndsAt:   c.EndsAt,
	})
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Competition, error) {
	row, err := r.db.Query.GetCompetitionByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Competition{}, ErrNotFound
		}
		return model.Competition{}, fmt.Errorf("get competition: %w", err)
	}
	return mapper.CompetitionFromDB(row), nil
}

func (r *PostgresRepository) JoinWithTradingAccount(
	ctx context.Context,
	competitionID uuid.UUID,
	userID uuid.UUID,
	login int64,
	broker string,
	investorPasswordEncrypted string,
) error {
	return r.db.WithTx(ctx, func(q *sqlc.Queries) error {
		// 1) Create trading account (try once)
		_, err := q.CreateTradingAccount(ctx, sqlc.CreateTradingAccountParams{
			Login:                     login,
			UserID:                    userID,
			Broker:                    broker,
			InvestorPasswordEncrypted: investorPasswordEncrypted,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				switch pgErr.ConstraintName {
				case "trading_accounts_pkey":
					return ErrLoginTaken
				case "trading_accounts_user_id_unique":
					return ErrAccountAlreadyExists
				default:
					return ErrAccountAlreadyExists
				}
			}
			return err
		}

		// 2) Join competition (ovo ti je OK kao :one)
		_, err = q.JoinCompetitionBeforeStart(ctx, sqlc.JoinCompetitionBeforeStartParams{
			CompetitionID:       competitionID,
			TradingAccountLogin: login,
		})
		if err == nil {
			return nil
		}

		// Already joined => unique violation on competition_members PK
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyJoined
		}

		// No rows => either not found OR already started (jer WHERE uslov nije prošao)
		if errors.Is(err, sql.ErrNoRows) {
			_, err2 := q.GetCompetitionStartTime(ctx, competitionID)
			if err2 == nil {
				return ErrAlreadyStarted
			}
			if errors.Is(err2, sql.ErrNoRows) {
				return ErrNotFound
			}
			return err2
		}

		return err
	})
}

func (r *PostgresRepository) UpdateAccountSize(ctx context.Context, competitionID uuid.UUID, login int64, accountSize float64) error {
	err := r.db.Query.UpdateCompetitionMemberAccountSize(ctx, sqlc.UpdateCompetitionMemberAccountSizeParams{
		CompetitionID:       competitionID,
		TradingAccountLogin: login,
		AccountSize:         accountSize,
	})
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrMemberNotFound
	}
	return fmt.Errorf("update account size: %w", err)
}

func (r *PostgresRepository) GetMemberAccountSize(ctx context.Context, competitionID uuid.UUID, login int64) (float64, error) {
	size, err := r.db.Query.GetCompetitionMemberAccountSize(ctx, sqlc.GetCompetitionMemberAccountSizeParams{
		CompetitionID:       competitionID,
		TradingAccountLogin: login,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotMember
		}
		return 0, fmt.Errorf("get account size: %w", err)
	}
	return size, nil
}

func (r *PostgresRepository) GetLeaderboard(ctx context.Context, competitionID uuid.UUID, limit, offset int32) ([]model.LeaderboardEntry, error) {
	rows, err := r.db.Query.GetCompetitionLeaderboard(ctx, sqlc.GetCompetitionLeaderboardParams{
		CompetitionID: competitionID,
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		return nil, fmt.Errorf("get leaderboard: %w", err)
	}
	return mapper.LeaderboardFromDB(rows), nil
}

func (r *PostgresRepository) InsertTrades(ctx context.Context, competitionID uuid.UUID, login int64, trades []model.Trade) error {
	for _, t := range trades {
		err := r.db.Query.InsertTrade(ctx, sqlc.InsertTradeParams{
			TradingAccountLogin: login,
			CompetitionID:       competitionID,
			PositionID:          t.PositionID,
			Symbol:              t.Symbol,
			Side:                t.Side,
			Volume:              t.Volume,
			OpenTime:            t.OpenTime,
			CloseTime:           t.CloseTime,
			OpenPrice:           t.OpenPrice,
			ClosePrice:          t.ClosePrice,
			Profit:              t.Profit,
			Commission:          t.Commission,
			Swap:                t.Swap,
		})
		if err == nil {
			continue
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			switch pgErr.ConstraintName {
			case "trades_competition_id_fkey":
				return ErrNotFound
			case "trades_trading_account_login_fkey":
				return ErrTradingAccountNotFound
			case "trades_member_fkey":
				return ErrNotMember
			}
		}
		return fmt.Errorf("insert trade (position_id=%d): %w", t.PositionID, err)
	}
	return nil
}

func (r *PostgresRepository) GetUserCompetitionState(ctx context.Context, userID, competitionID uuid.UUID) (sqlc.GetCompetitionUserStateRow, error) {
	state, err := r.db.Query.GetCompetitionUserState(ctx, sqlc.GetCompetitionUserStateParams{
		UserID:        userID,
		CompetitionID: competitionID,
	})
	if err != nil {
		return sqlc.GetCompetitionUserStateRow{}, fmt.Errorf("get user state: %w", err)
	}
	return state, nil
}

func (r *PostgresRepository) GetCurrent(ctx context.Context) (sqlc.Competition, error) {
	comp, err := r.db.Query.GetCurrentCompetition(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Competition{}, ErrNotFound
		}
		return sqlc.Competition{}, fmt.Errorf("get current competition: %w", err)
	}
	return comp, nil
}

func (r *PostgresRepository) CreateAccountRequest(ctx context.Context, userID, competitionID uuid.UUID) error {
	err := r.db.Query.CreateCompetitionAccountRequest(ctx, sqlc.CreateCompetitionAccountRequestParams{
		UserID:        userID,
		CompetitionID: competitionID,
	})
	if err != nil {
		return fmt.Errorf("create account request: %w", err)
	}
	return nil
}
