package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/filipcvejic/trading_tournament/db"
	"github.com/filipcvejic/trading_tournament/db/sqlc"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, email, username, discordUsername, passwordHash string) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	UpdatePasswordHash(ctx context.Context, userID uuid.UUID, hash string) error
}

type PostgresRepository struct {
	db *db.DB
}

func NewPostgresRepository(database *db.DB) *PostgresRepository {
	return &PostgresRepository{db: database}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	email, username, discordUsername, passwordHash string,
) (User, error) {
	row, err := r.db.Query.CreateUser(ctx, sqlc.CreateUserParams{
		Email:           email,
		Username:        username,
		DiscordUsername: discordUsername,
		PasswordHash:    passwordHash,
	})
	if err != nil {
		return User{}, err
	}

	return User{
		ID:              row.ID,
		Email:           row.Email,
		Username:        row.Username,
		DiscordUsername: row.DiscordUsername,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	row, err := r.db.Query.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return User{
		ID:              row.ID,
		Email:           row.Email,
		Username:        row.Username,
		DiscordUsername: row.DiscordUsername,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	row, err := r.db.Query.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, err
	}

	return User{
		ID:              row.ID,
		Email:           row.Email,
		Username:        row.Username,
		DiscordUsername: row.DiscordUsername,
		PasswordHash:    row.PasswordHash,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}, nil
}

func (r *PostgresRepository) UpdatePasswordHash(ctx context.Context, userID uuid.UUID, hash string) error {
	return r.db.Query.UpdatePasswordHash(ctx, sqlc.UpdatePasswordHashParams{
		ID:           userID,
		PasswordHash: hash,
	})
}
