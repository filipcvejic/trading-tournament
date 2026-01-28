package http

import (
	"errors"
	"net/http"

	"github.com/filipcvejic/trading_tournament/internal/auth"
	"github.com/filipcvejic/trading_tournament/internal/competition"
	"github.com/filipcvejic/trading_tournament/internal/httputil"
)

type errorMapping struct {
	status  int
	message string
}

var errorMap = map[error]errorMapping{
	// Not Found (404)
	competition.ErrNotFound:               {http.StatusNotFound, "Competition not found"},
	competition.ErrMemberNotFound:         {http.StatusNotFound, "Competition member not found"},
	competition.ErrTradingAccountNotFound: {http.StatusNotFound, "Trading account not found"},

	// Conflict (409)
	competition.ErrAlreadyStarted:       {http.StatusConflict, "Competition has already started"},
	competition.ErrAlreadyJoined:        {http.StatusConflict, "You have already joined this competition"},
	competition.ErrAccountAlreadyExists: {http.StatusConflict, "You already have a trading account"},
	competition.ErrLoginTaken:           {http.StatusConflict, "This trading account login is already taken"},

	// Forbidden (403)
	competition.ErrNotMember: {http.StatusForbidden, "You are not a member of this competition"},

	// Bad Request (400)
	competition.ErrInvalidName:             {http.StatusBadRequest, "Competition name cannot be empty"},
	competition.ErrInvalidTimeRange:        {http.StatusBadRequest, "End time must be after start time"},
	competition.ErrInvalidAccountSize:      {http.StatusBadRequest, "Account size must be greater than zero"},
	competition.ErrAccountSizeNotSet:       {http.StatusBadRequest, "Account size must be set before inserting trades"},
	competition.ErrInvalidLogin:            {http.StatusBadRequest, "Trading account login must be a positive number"},
	competition.ErrInvalidPositionID:       {http.StatusBadRequest, "Position ID must be a positive number"},
	competition.ErrInvalidSymbol:           {http.StatusBadRequest, "Symbol cannot be empty"},
	competition.ErrInvalidSide:             {http.StatusBadRequest, "Side must be 'buy' or 'sell'"},
	competition.ErrInvalidTradeTimeRange:   {http.StatusBadRequest, "Trade close time must be after open time"},
	competition.ErrInvalidBroker:           {http.StatusBadRequest, "Broker cannot be empty"},
	competition.ErrInvalidInvestorPassword: {http.StatusBadRequest, "Investor password cannot be empty"},

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
