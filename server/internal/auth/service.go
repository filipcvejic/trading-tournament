package auth

import (
	"context"
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type AuthService struct {
	userRepo user.Repository
	//refreshTokenRepo RefreshTokenRepository
	jwtSecret      []byte
	accessTokenTTL time.Duration
}

func NewAuthService(userRepo user.Repository, refreshTokenRepo RefreshTokenRepository, jwtSecret string, accessTokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		//refreshTokenRepo: refreshTokenRepo,
		jwtSecret: []byte(jwtSecret),
		//accessTokenTTL:   accessTokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, discordUsername, password string) (*user.User, error) {
	if email == "" || username == "" || discordUsername == "" || password == "" {
		return nil, ErrInvalidInput
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.Create(ctx, email, username, discordUsername, hashedPassword)
	if err == nil {
		return &user, nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_email_unique":
			return nil, ErrEmailInUse
		case "users_username_unique":
			return nil, ErrUsernameInUse
		case "users_discord_username_unique":
			return nil, ErrDiscordUsernameInUse
		}
	}

	return nil, err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", ErrInvalidInput
	}

	user, err := s.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.generateAccessToken(user)
}

//func (s *AuthService) LoginWithRefresh(ctx context.Context, email, password string, refreshTokenTTL time.Duration) (accessToken string, refreshToken string, err error) {
//	if email == "" || password == "" {
//		return "", "", ErrInvalidInput
//	}
//
//	user, err := s.userRepo.GetByEmail(ctx, email)
//	if err != nil {
//		return "", "", ErrInvalidCredentials
//	}
//
//	if err := VerifyPassword(user.PasswordHash, password); err != nil {
//		return "", "", ErrInvalidCredentials
//	}
//
//	accessToken, err = s.generateAccessToken(user)
//	if err != nil {
//		return "", "", err
//	}
//
//	token, err := s.refreshTokenRepo.Create(ctx, user.ID, refreshTokenTTL)
//	if err != nil {
//		return "", "", err
//	}
//
//	return accessToken, token.Token, nil
//}
//
//func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenString string) (string, error) {
//	if refreshTokenString == "" {
//		return "", ErrInvalidInput
//	}
//
//	rt, err := s.refreshTokenRepo.Get(ctx, refreshTokenString)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return "", ErrInvalidToken
//		}
//		return "", err
//	}
//
//	if rt.Revoked {
//		return "", ErrInvalidToken
//	}
//
//	if time.Now().After(rt.ExpiresAt) {
//		return "", ErrExpiredToken
//	}
//
//	user, err := s.userRepo.GetByID(ctx, rt.UserID)
//	if err != nil {
//		return "", err
//	}
//
//	return s.generateAccessToken(user)
//}

func (s *AuthService) Me(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, ErrUnauthorized
	}

	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
