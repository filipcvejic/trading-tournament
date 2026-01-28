package http

import (
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/auth"
	"github.com/filipcvejic/trading_tournament/internal/httputil"
	"github.com/filipcvejic/trading_tournament/internal/user"
	"net/http"
)

type errorMapping struct {
	status  int
	message string
}

var errorMap = map[error]errorMapping{
	// Not Found (404)
	user.ErrNotFound: {http.StatusNotFound, "User not found"},

	// Conflict (409)
	user.ErrUsernameAlreadyExists:        {http.StatusConflict, "Username is already taken"},
	user.ErrDiscordUsernameAlreadyExists: {http.StatusConflict, "Discord username is already taken"},
	user.ErrEmailAlreadyExists:           {http.StatusConflict, "Email is already in use"},

	// Bad Request (400)
	user.ErrInvalidEmail: {
		http.StatusBadRequest,
		"Email is invalid",
	},
	user.ErrInvalidUsername: {
		http.StatusBadRequest,
		"Username must be between 3 and 20 characters and must not contain whitespace",
	},
	user.ErrInvalidDiscordUsername: {
		http.StatusBadRequest,
		"Discord username must be at least 2 characters and must not contain whitespace",
	},
	user.ErrInvalidPassword: {
		http.StatusBadRequest,
		"Password must be at least 11 characters and include 1 lowercase, 1 uppercase, and 1 special character",
	},
	auth.ErrInvalidInput: {
		http.StatusBadRequest, "Invalid input",
	},
	auth.ErrInvalidCredentials: {http.StatusBadRequest, "Invalid email or password"},
	
	auth.ErrUnauthorized: {http.StatusUnauthorized, "Unauthorized"},
}

// writeDomainError maps domain errors to HTTP responses
func writeDomainError(w http.ResponseWriter, r *http.Request, err error) {
	for domainErr, mapping := range errorMap {
		if errors.Is(err, domainErr) {
			httputil.WriteError(w, r, mapping.status, mapping.message, err)
			return
		}
	}

	// Unknown error
	httputil.WriteInternalError(w, r, err)
}
