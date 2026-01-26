package competition

import (
	"context"
	"database/sql"
	"errors"
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

func (r *PostgresRepository) Create(
	ctx context.Context,
	c model.Competition,
) error {
	_, err := r.db.Query.CreateCompetition(ctx, sqlc.CreateCompetitionParams{
		ID:       c.ID,
		Name:     c.Name,
		StartsAt: c.StartsAt,
		EndsAt:   c.EndsAt,
	})

	return err
}

func (r *PostgresRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (model.Competition, error) {
	row, err := r.db.Query.GetCompetitionByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Competition{}, ErrCompetitionNotFound
		}
		return model.Competition{}, err
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
		// 1) Ako user već ima trading account (user_id UNIQUE), ne dozvoli drugi login.
		if existing, err := q.GetTradingAccountByUserID(ctx, userID); err == nil {
			if existing.Login != login {
				return ErrTradingAccountAlreadyExistsForUser
			}
			// isti login -> ok (idempotent)
		} else if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		// 2) Create trading account (ako već postoji za tog usera/login, tretiramo kao idempotent)
		_, err := q.CreateTradingAccount(ctx, sqlc.CreateTradingAccountParams{
			Login:                     login,
			UserID:                    userID,
			Broker:                    broker,
			InvestorPasswordEncrypted: investorPasswordEncrypted,
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				// duplikat: login PK ili user_id unique
				// proveri ko je vlasnik login-a
				acc, accErr := q.GetTradingAccountByLogin(ctx, login)
				if accErr == nil {
					if acc.UserID != userID {
						return ErrTradingAccountLoginTaken
					}
					// isti user -> idempotent
				} else if errors.Is(accErr, sql.ErrNoRows) {
					// nema po login-u, znači user_id unique je udario (drugi login)
					return ErrTradingAccountAlreadyExistsForUser
				} else {
					return accErr
				}
			} else {
				return err
			}
		}

		// 3) Join competition pre starta
		_, err = q.JoinCompetitionBeforeStart(ctx, sqlc.JoinCompetitionBeforeStartParams{
			CompetitionID:       competitionID,
			TradingAccountLogin: login,
		})
		if err == nil {
			return nil
		}

		// već joined
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyJoined
		}

		// nije insertovalo ništa -> (ne postoji competition) ili (već počelo)
		if errors.Is(err, sql.ErrNoRows) {
			_, err2 := q.GetCompetitionStartTime(ctx, competitionID)
			if err2 == nil {
				return ErrCompetitionAlreadyStarted
			}
			if errors.Is(err2, sql.ErrNoRows) {
				return ErrCompetitionNotFound
			}
			return err2
		}

		return err
	})
}

func (r *PostgresRepository) UpdateAccountSize(ctx context.Context, competitionID uuid.UUID, login int64, accountSize float64) error {
	if accountSize <= 0 {
		return ErrInvalidAccountSize
	}

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

	return err
}

func (r *PostgresRepository) GetMemberAccountSize(
	ctx context.Context,
	competitionID uuid.UUID,
	login int64,
) (float64, error) {

	size, err := r.db.Query.GetCompetitionMemberAccountSize(ctx, sqlc.GetCompetitionMemberAccountSizeParams{
		CompetitionID:       competitionID,
		TradingAccountLogin: login,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotCompetitionMember
		}
		return 0, err
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
		return nil, err
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
				return ErrCompetitionNotFound
			case "trades_trading_account_login_fkey":
				return ErrTradingAccountNotFound
			case "trades_member_fkey":
				return ErrNotCompetitionMember
			}
		}
		return err
	}

	return nil
}

func (r *PostgresRepository) GetUserCompetitionState(ctx context.Context, userID, competitionID uuid.UUID) (sqlc.GetCompetitionUserStateRow, error) {
	return r.db.Query.GetCompetitionUserState(ctx, sqlc.GetCompetitionUserStateParams{
		UserID:        userID,
		CompetitionID: competitionID,
	})
}

func (r *PostgresRepository) GetCurrent(ctx context.Context) (sqlc.Competition, error) {
	return r.db.Query.GetCurrentCompetition(ctx)
}

func (r *PostgresRepository) CreateAccountRequest(ctx context.Context, userID, competitionID uuid.UUID) error {
	return r.db.Query.CreateCompetitionAccountRequest(ctx, sqlc.CreateCompetitionAccountRequestParams{
		UserID:        userID,
		CompetitionID: competitionID,
	})
}
