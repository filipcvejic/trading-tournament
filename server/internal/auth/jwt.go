package auth

import (
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (s *AuthService) generateAccessToken(user user.User) (string, error) {
	now := time.Now()
	expirationTime := now.Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"exp":   expirationTime.Unix(),
		"iat":   now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
