package competition

import (
	"context"
	"encoding/base64"
	"github.com/filipcvejic/trading_tournament/internal/auth"
	"github.com/filipcvejic/trading_tournament/internal/competition/dto"
	"github.com/filipcvejic/trading_tournament/internal/competition/model"
	"github.com/filipcvejic/trading_tournament/internal/crypto"
	"github.com/google/uuid"
	"strings"
)

type Service struct {
	repo      Repository
	cryptoKey []byte
}

func NewService(repo Repository, cryptoKeyBase64 string) (*Service, error) {
	key, err := base64.StdEncoding.DecodeString(cryptoKeyBase64)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, crypto.ErrInvalidKeyLength
	}
	return &Service{repo: repo, cryptoKey: key}, nil
}

func (s *Service) Create(
	ctx context.Context,
	c model.Competition,
) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return ErrInvalidCompetitionName
	}

	if !c.EndsAt.After(c.StartsAt) {
		return ErrInvalidCompetitionTimeRange
	}

	return s.repo.Create(ctx, c)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (model.Competition, error) {
	if id == uuid.Nil {
		return model.Competition{}, ErrCompetitionNotFound
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) JoinWithTradingAccount(
	ctx context.Context,
	competitionID, userID uuid.UUID,
	login int64,
	broker string,
	investorPassword string,
) error {
	if competitionID == uuid.Nil {
		return ErrCompetitionNotFound
	}
	if userID == uuid.Nil {
		return auth.ErrUnauthorized
	}
	if login <= 0 {
		return ErrInvalidTradingAccountLogin
	}
	if broker == "" {
		return ErrInvalidBroker
	}
	if investorPassword == "" {
		return ErrInvalidInvestorPassword
	}

	encrypted, err := crypto.EncryptString(s.cryptoKey, investorPassword)
	if err != nil {
		return err
	}

	return s.repo.JoinWithTradingAccount(ctx, competitionID, userID, login, broker, encrypted)
}

func (s *Service) UpdateAccountSize(
	ctx context.Context,
	competitionID uuid.UUID,
	login int64,
	accountSize float64,
) error {
	if competitionID == uuid.Nil {
		return ErrCompetitionNotFound
	}
	if login <= 0 {
		return ErrInvalidTradingAccountLogin
	}
	if accountSize <= 0 {
		return ErrInvalidAccountSize
	}

	return s.repo.UpdateAccountSize(ctx, competitionID, login, accountSize)
}

const (
	defaultLeaderboardLimit int32 = 50
	maxLeaderboardLimit     int32 = 200
)

func (s *Service) GetLeaderboard(
	ctx context.Context,
	competitionID uuid.UUID,
	limit int32,
	offset int32,
) ([]model.LeaderboardEntry, error) {
	if limit <= 0 {
		limit = defaultLeaderboardLimit
	}
	if limit > maxLeaderboardLimit {
		limit = maxLeaderboardLimit
	}
	if offset < 0 {
		offset = 0
	}

	entries, err := s.repo.GetLeaderboard(ctx, competitionID, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		_, err := s.repo.GetByID(ctx, competitionID)
		if err != nil {
			return nil, err
		}
	}

	return entries, nil
}

func (s *Service) InsertTrades(ctx context.Context, competitionID uuid.UUID, login int64, trades []model.Trade) error {
	if competitionID == uuid.Nil {
		return ErrCompetitionNotFound
	}
	if login <= 0 {
		return ErrInvalidTradingAccountLogin
	}

	size, err := s.repo.GetMemberAccountSize(ctx, competitionID, login)
	if err != nil {
		return err
	}
	if size == 0 {
		return ErrAccountSizeNotSet
	}

	for i := range trades {
		t := trades[i]
		if t.PositionID <= 0 {
			return ErrInvalidPositionID
		}
		if t.Symbol == "" {
			return ErrInvalidSymbol
		}
		if t.Side == "" {
			return ErrInvalidSide
		}
		if !t.CloseTime.After(t.OpenTime) {
			return ErrInvalidTradeTimeRange
		}
	}

	return s.repo.InsertTrades(ctx, competitionID, login, trades)
}

func (s *Service) GetUserCompetitionState(
	ctx context.Context,
	userID,
	competitionID uuid.UUID,
) (*dto.CompetitionUserStateResponse, error) {

	state, err := s.repo.GetUserCompetitionState(ctx, userID, competitionID)
	if err != nil {
		return nil, err
	}

	return &dto.CompetitionUserStateResponse{
		HasRequestedAccount: state.HasRequestedAccount,
		HasJoined:           state.HasJoined,
	}, nil
}

func (s *Service) GetCurrentCompetition(ctx context.Context) (*dto.CompetitionResponse, error) {
	c, err := s.repo.GetCurrent(ctx)
	if err != nil {
		return nil, ErrCompetitionNotFound
	}

	return &dto.CompetitionResponse{
		ID:       c.ID,
		Name:     c.Name,
		StartsAt: c.StartsAt,
		EndsAt:   c.EndsAt,
	}, nil
}

func (s *Service) RequestAccount(ctx context.Context, userID, competitionID uuid.UUID) error {
	return s.repo.CreateAccountRequest(ctx, userID, competitionID)
}
