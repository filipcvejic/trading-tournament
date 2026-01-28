package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/filipcvejic/trading_tournament/internal/auth"
	"github.com/filipcvejic/trading_tournament/internal/competition"
	"github.com/filipcvejic/trading_tournament/internal/competition/dto"
	"github.com/filipcvejic/trading_tournament/internal/competition/mapper"
	"github.com/filipcvejic/trading_tournament/internal/competition/model"
	"github.com/filipcvejic/trading_tournament/internal/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *competition.Service
}

func NewHandler(service *competition.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/competitions", func(r chi.Router) {
		r.Post("/", h.createCompetition)
		r.Post("/{competitionID}/members/{accountLogin}/account-size", h.updateAccountSize)
		r.Post("/{competitionID}/trades", h.insertTrades)

		r.Group(func(r chi.Router) {
			r.Use(auth.AuthenticationMiddleware)
			r.Get("/{competitionID}", h.getCompetitionByID)
			r.Get("/current", h.getCurrent)
			r.Get("/{competitionID}/leaderboard", h.getLeaderboard)
			r.Post("/{competitionID}/join", h.joinCompetition)
			r.Get("/{competitionID}/me", h.getMe)
			r.Post("/{competitionID}/account-requests", h.requestAccount)
		})
	})
}

func (h *Handler) createCompetition(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCompetitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteClientError(w, r, "Invalid JSON body", err)
		return
	}

	c := model.Competition{
		ID:       uuid.New(),
		Name:     req.Name,
		StartsAt: req.StartsAt,
		EndsAt:   req.EndsAt,
	}

	if err := h.service.Create(r.Context(), c); err != nil {
		writeDomainError(w, r, err)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, dto.CompetitionResponse{
		ID:       c.ID,
		Name:     c.Name,
		StartsAt: c.StartsAt,
		EndsAt:   c.EndsAt,
	})
}

func (h *Handler) getCompetitionByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "competitionID"))
	if err != nil {
		httputil.WriteClientError(w, r, "Invalid competition ID format", err)
		return
	}

	c, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, dto.CompetitionResponse{
		ID:       c.ID,
		Name:     c.Name,
		StartsAt: c.StartsAt,
		EndsAt:   c.EndsAt,
	})
}

func (h *Handler) joinCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := uuid.Parse(chi.URLParam(r, "competitionID"))
	if err != nil {
		httputil.WriteClientError(w, r, "Invalid competition ID format", err)
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		httputil.WriteUnauthorized(w, r)
		return
	}

	var req dto.JoinCompetitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteClientError(w, r, "Invalid JSON body", err)
		return
	}

	if err := h.service.JoinWithTradingAccount(r.Context(), competitionID, userID, req.Login, req.Broker, req.InvestorPassword); err != nil {
		writeDomainError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateAccountSize(w http.ResponseWriter, r *http.Request) {
	competitionID, err := uuid.Parse(chi.URLParam(r, "competitionID"))
	if err != nil {
		httputil.WriteClientError(w, r, "Invalid competition ID format", err)
		return
	}

	accountLogin, err := strconv.ParseInt(chi.URLParam(r, "accountLogin"), 10, 64)
	if err != nil {
		httputil.WriteClientError(w, r, "Invalid account login format", err)
		return
	}

	var req dto.UpdateAccountSizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteClientError(w, r, "Invalid JSON body", err)
		return
	}

	if err := h.service.UpdateAccountSize(r.Context(), competitionID, accountLogin, req.AccountSize); err != nil {
		writeDomainError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getLeaderboard(w http.ResponseWriter, r *http.Request) {
	competitionID, err := uuid.Parse(chi.URLParam(r, "competitionID"))
	if err != nil {
		httputil.WriteClientError(w, r, "Invalid competition ID format", err)
		return
	}

	var limit, offset int32

	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			httputil.WriteClientError(w, r, "Invalid limit parameter", err)
			return
		}
		limit = int32(n)
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			httputil.WriteClientError(w, r, "Invalid offset parameter", err)
			return
		}
		offset = int32(n)
	}

	entries, err := h.service.GetLeaderboard(r.Context(), competitionID, limit, offset)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, mapper.LeaderboardToDTO(entries))
}

func (h *Handler) insertTrades(w http.ResponseWriter, r *http.Request) {
	competitionID, err := uuid.Parse(chi.URLParam(r, "competitionID"))
	if err != nil {
		httputil.WriteClientError(w, r, "Invalid competition ID format", err)
		return
	}

	var req dto.InsertTradesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteClientError(w, r, "Invalid JSON body", err)
		return
	}

	trades := make([]model.Trade, len(req.Trades))
	for i, it := range req.Trades {
		trades[i] = model.Trade{
			PositionID: it.PositionID,
			Symbol:     it.Symbol,
			Side:       it.Side,
			Volume:     it.Volume,
			OpenTime:   it.OpenTime,
			CloseTime:  it.CloseTime,
			OpenPrice:  it.OpenPrice,
			ClosePrice: it.ClosePrice,
			Profit:     it.Profit,
			Commission: it.Commission,
			Swap:       it.Swap,
		}
	}

	if err := h.service.InsertTrades(r.Context(), competitionID, req.TradingAccountLogin, trades); err != nil {
		writeDomainError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	competitionID, err := uuid.Parse(chi.URLParam(r, "competitionID"))
	if err != nil {
		httputil.WriteClientError(w, r, "Invalid competition ID format", err)
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		httputil.WriteUnauthorized(w, r)
		return
	}

	state, err := h.service.GetUserCompetitionState(r.Context(), userID, competitionID)
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, state)
}

func (h *Handler) getCurrent(w http.ResponseWriter, r *http.Request) {
	comp, err := h.service.GetCurrentCompetition(r.Context())
	if err != nil {
		writeDomainError(w, r, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, comp)
}

func (h *Handler) requestAccount(w http.ResponseWriter, r *http.Request) {
	competitionID, err := uuid.Parse(chi.URLParam(r, "competitionID"))
	if err != nil {
		httputil.WriteClientError(w, r, "Invalid competition ID format", err)
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		httputil.WriteUnauthorized(w, r)
		return
	}

	if err := h.service.RequestAccount(r.Context(), userID, competitionID); err != nil {
		writeDomainError(w, r, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
