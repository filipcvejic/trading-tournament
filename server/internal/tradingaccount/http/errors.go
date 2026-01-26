package http

import (
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/tradingaccount"
	"net/http"
)

func writeTradingAccountError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, tradingaccount.ErrTradingAccountNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, tradingaccount.ErrLoginAlreadyExists):
		http.Error(w, err.Error(), http.StatusConflict)

	case errors.Is(err, tradingaccount.ErrUserNotFound):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, tradingaccount.ErrInvalidLogin),
		errors.Is(err, tradingaccount.ErrInvalidUserID),
		errors.Is(err, tradingaccount.ErrInvalidBroker),
		errors.Is(err, tradingaccount.ErrInvalidInvestorPassword):
		http.Error(w, err.Error(), http.StatusBadRequest)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
