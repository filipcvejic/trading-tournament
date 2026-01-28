package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, email, username, discordUsername, password string) (User, error) {
	if email == "" || strings.TrimSpace(email) != email {
		return User{}, ErrInvalidEmail
	}
	if username == "" {
		return User{}, ErrInvalidUsername
	}
	if discordUsername == "" {
		return User{}, ErrInvalidDiscordUsername
	}
	if password == "" {
		return User{}, ErrInvalidPassword
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	passwordHash := string(hashBytes)

	u, err := s.repo.Create(ctx, email, username, discordUsername, passwordHash)
	if err == nil {
		return u, nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_email_unique":
			return User{}, ErrEmailAlreadyExists
		case "users_username_unique":
			return User{}, ErrUsernameAlreadyExists
		case "users_discord_username_unique":
			return User{}, ErrDiscordUsernameAlreadyExists
		}
	}

	return User{}, err
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (User, error) {
	if id == uuid.Nil {
		return User{}, ErrNotFound
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (User, error) {
	if email == "" {
		return User{}, ErrNotFound
	}
	return s.repo.GetByEmail(ctx, email)
}
