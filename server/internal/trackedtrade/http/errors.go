package http

import (
	"errors"
	"github.com/filipcvejic/trading_tournament/internal/httputil"
	"github.com/filipcvejic/trading_tournament/internal/trackedtrade"
	nethttp "net/http"
)

func writeDomainError(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
	switch {
	case errors.Is(err, trackedtrade.ErrInvalidEventType):
		httputil.WriteClientError(w, r, "Invalid event type", err)

	case errors.Is(err, trackedtrade.ErrMissingOpenFields):
		httputil.WriteClientError(w, r, "Missing required fields for OPEN event", err)

	case errors.Is(err, trackedtrade.ErrMissingCloseFields):
		httputil.WriteClientError(w, r, "Missing required fields for CLOSE event", err)

	case errors.Is(err, trackedtrade.ErrMissingUpdateFields):
		httputil.WriteClientError(w, r, "Missing required fields for UPDATE event", err)

	case errors.Is(err, trackedtrade.ErrTradeNotFound):
		httputil.WriteError(w, r, nethttp.StatusNotFound, "Tracked trade not found", err)

	default:
		httputil.WriteInternalError(w, r, err)
	}
}
