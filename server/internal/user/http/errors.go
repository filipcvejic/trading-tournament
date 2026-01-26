package http

import (
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/user"
	"net/http"
)

func writeUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, user.ErrEmailAlreadyExists),
		errors.Is(err, user.ErrUsernameAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)
	
	case errors.Is(err, user.ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrInvalidUsername),
		errors.Is(err, user.ErrInvalidPassword):
		http.Error(w, err.Error(), http.StatusBadRequest)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
