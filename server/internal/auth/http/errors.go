package http

import (
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/auth"
	"net/http"
)

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, auth.ErrEmailInUse),
		errors.Is(err, auth.ErrUsernameInUse),
		errors.Is(err, auth.ErrDiscordUsernameInUse):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, auth.ErrInvalidCredentials):
		http.Error(w, "invalid credentials", http.StatusUnauthorized)

	case errors.Is(err, auth.ErrInvalidToken),
		errors.Is(err, auth.ErrExpiredToken):
		http.Error(w, err.Error(), http.StatusUnauthorized)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
