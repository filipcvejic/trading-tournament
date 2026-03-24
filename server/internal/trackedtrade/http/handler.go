package http

import (
	"encoding/json"
	"github.com/filipcvejic/trading_tournament/internal/auth"
	"github.com/go-chi/chi/v5"
	nethttp "net/http"

	"github.com/filipcvejic/trading_tournament/internal/httputil"
	"github.com/filipcvejic/trading_tournament/internal/trackedtrade"
	"github.com/filipcvejic/trading_tournament/internal/validation"
)

type Handler struct {
	service *trackedtrade.Service
}

func NewHandler(service *trackedtrade.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		r.Use(auth.AuthenticationMiddleware)
		r.Use(auth.RequireAdmin)

		r.Get("/tracked-trades", h.List)
	})

	r.Post("/tracked-trades/events", h.IngestEvent)
}

func (h *Handler) IngestEvent(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req trackedtrade.IngestTrackedTradeEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteClientError(w, r, "Neispravan JSON body", err)
		return
	}

	if err := validation.V.Struct(req); err != nil {
		httputil.WriteClientError(w, r, validation.FirstMessage(err), err)
		return
	}

	if err := h.service.IngestEvent(r.Context(), req); err != nil {
		writeDomainError(w, r, err)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *Handler) List(w nethttp.ResponseWriter, r *nethttp.Request) {
	trades, err := h.service.List(r.Context())
	if err != nil {
		httputil.WriteInternalError(w, r, err)
		return
	}

	response := make([]trackedtrade.TrackedTradeResponse, 0, len(trades))

	for _, t := range trades {
		response = append(response, trackedtrade.TrackedTradeResponse{
			PositionID: t.PositionID,
			Symbol:     t.Symbol,
			Side:       string(t.Side),
			OpenPrice:  t.OpenPrice,
			Volume:     t.Volume,
			StopLoss:   t.StopLoss,
			OpenedAt:   t.OpenedAt,
			ClosedAt:   t.ClosedAt,
		})
	}

	httputil.WriteJSON(w, nethttp.StatusOK, response)
}
