package http

import (
	"encoding/json"
	"github.com/filipcvejic/trading_tournament/internal/httputil"
	"github.com/filipcvejic/trading_tournament/internal/tradingaccount"
	"github.com/filipcvejic/trading_tournament/internal/validation"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Handler struct {
	service *tradingaccount.Service
}

func NewHandler(service *tradingaccount.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/trading-accounts", func(r chi.Router) {
		r.Post("/", h.createTradingAccount)
		r.Get("/{login}", h.getByLogin)
		r.Get("/{login}/trade-history", h.getTradeHistory)
	})
}

func (h *Handler) createTradingAccount(w http.ResponseWriter, r *http.Request) {
	var req tradingaccount.CreateTradingAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteClientError(w, r, "Invalid JSON body", err)
		return
	}

	if err := validation.V.Struct(req); err != nil {
		httputil.WriteClientError(w, r, validation.FirstMessage(err), err)
		return
	}

	acc, err := h.service.Create(
		r.Context(),
		req.Login,
		req.UserID,
		req.Broker,
		req.InvestorPassword,
	)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, tradingaccount.TradingAccountResponse{
		Login:     acc.Login,
		UserID:    acc.UserID,
		Broker:    acc.Broker,
		CreatedAt: acc.CreatedAt,
	})
}

func (h *Handler) getByLogin(w http.ResponseWriter, r *http.Request) {
	loginStr := chi.URLParam(r, "login")
	login, err := strconv.ParseInt(loginStr, 10, 64)
	if err != nil || login <= 0 {
		httputil.WriteClientError(w, r, "Invalid login format", err)
		return
	}

	acc, err := h.service.GetByLogin(r.Context(), login)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, tradingaccount.TradingAccountResponse{
		Login:     acc.Login,
		UserID:    acc.UserID,
		Broker:    acc.Broker,
		CreatedAt: acc.CreatedAt,
	})
}

func (h *Handler) getTradeHistory(w http.ResponseWriter, r *http.Request) {
	loginStr := chi.URLParam(r, "login")
	login, err := strconv.ParseInt(loginStr, 10, 64)
	if err != nil || login <= 0 {
		httputil.WriteClientError(w, r, "Invalid login format", err)
		return
	}

	resp, err := h.service.GetTradeHistory(r.Context(), login)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, resp)
}
