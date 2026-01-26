package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/filipcvejic/trading_tournament/db"
	"github.com/filipcvejic/trading_tournament/db/sqlc"
	"github.com/filipcvejic/trading_tournament/internal/auth/model"
	"github.com/google/uuid"
	"time"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (model.RefreshToken, error)
	Get(ctx context.Context, token string) (model.RefreshToken, error)
	Revoke(ctx context.Context, token string) error
}

type PostgresRefreshTokenRepository struct {
	db *db.DB
}

func NewPostgresRefreshTokenRepository(database *db.DB) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: database}
}

func (r *PostgresRefreshTokenRepository) Create(ctx context.Context, userID uuid.UUID, ttl time.Duration) (model.RefreshToken, error) {
	token, err := generateRefreshTokenString(32)
	if err != nil {
		return model.RefreshToken{}, err
	}

	expiresAt := time.Now().Add(ttl)

	row, err := r.db.Query.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return model.RefreshToken{}, err
	}

	return model.RefreshToken{
		Token:     row.Token,
		UserID:    row.UserID,
		ExpiresAt: row.ExpiresAt,
		Revoked:   row.Revoked,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *PostgresRefreshTokenRepository) Get(ctx context.Context, token string) (model.RefreshToken, error) {
	row, err := r.db.Query.GetRefreshToken(ctx, token)
	if err != nil {
		return model.RefreshToken{}, err
	}

	return model.RefreshToken{
		Token:     row.Token,
		UserID:    row.UserID,
		ExpiresAt: row.ExpiresAt,
		Revoked:   row.Revoked,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *PostgresRefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	return r.db.Query.RevokeRefreshToken(ctx, token)
}
func generateRefreshTokenString(nBytes int) (string, error) {
	if nBytes < 16 {
		return "", errors.New("refresh token length too small")
	}

	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// URL-safe, bez padding-a
	return base64.RawURLEncoding.EncodeToString(b), nil
}
