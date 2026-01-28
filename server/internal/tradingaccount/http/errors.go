package http

import (
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/tradingaccount"
	"github.com/filipcvejic/trading_tournament/internal/user"
	"net/http"

	"github.com/filipcvejic/trading_tournament/internal/auth"
	"github.com/filipcvejic/trading_tournament/internal/httputil"
)

type errorMapping struct {
	status  int
	message string
}

var errorMap = map[error]errorMapping{
	// Not Found (404)
	tradingaccount.ErrNotFound: {http.StatusNotFound, "Trading account not found"},
	user.ErrNotFound:           {http.StatusNotFound, "User not found"},

	// Conflict (409)
	tradingaccount.ErrLoginTaken: {
		http.StatusConflict,
		"Trading account with this login already exists",
	},

	// Bad Request (400)
	tradingaccount.ErrInvalidLogin: {
		http.StatusBadRequest,
		"Login must be at least 5 digits",
	},
	tradingaccount.ErrInvalidBroker: {
		http.StatusBadRequest,
		"Broker name must be at least 2 characters",
	},
	tradingaccount.ErrInvalidInvestorPassword: {
		http.StatusBadRequest,
		"Investor password must be at least 5 characters",
	},

	// Auth errors
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
