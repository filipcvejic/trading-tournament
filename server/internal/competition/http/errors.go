package http

import (
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/competition"
	"net/http"
)

func writeCompetitionError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, competition.ErrCompetitionNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, competition.ErrCompetitionAlreadyStarted),
		errors.Is(err, competition.ErrAlreadyJoined):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, competition.ErrNotCompetitionMember):
		http.Error(w, err.Error(), http.StatusForbidden)

	case errors.Is(err, competition.ErrInvalidCompetitionName),
		errors.Is(err, competition.ErrInvalidCompetitionTimeRange),
		errors.Is(err, competition.ErrInvalidTradingAccountLogin),
		errors.Is(err, competition.ErrInvalidAccountSize),
		errors.Is(err, competition.ErrInvalidPositionID),
		errors.Is(err, competition.ErrInvalidSymbol),
		errors.Is(err, competition.ErrInvalidSide),
		errors.Is(err, competition.ErrInvalidTradeTimeRange),
		errors.Is(err, competition.ErrAccountSizeNotSet):
		http.Error(w, err.Error(), http.StatusBadRequest)

	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
