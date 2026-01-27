package http

import (
	"encoding/json"
	"github.com/filipcvejic/trading_tournament/internal/auth"
	"github.com/filipcvejic/trading_tournament/internal/competition"
	"github.com/filipcvejic/trading_tournament/internal/competition/dto"
	"github.com/filipcvejic/trading_tournament/internal/competition/mapper"
	"github.com/filipcvejic/trading_tournament/internal/competition/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"strconv"
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

		// protected
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
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	c := model.Competition{
		ID:       uuid.New(),
		Name:     req.Name,
		StartsAt: req.StartsAt,
		EndsAt:   req.EndsAt,
	}

	if err := h.service.Create(r.Context(), c); err != nil {
		writeCompetitionError(w, err)
		return
	}

	resp := dto.CompetitionResponse{
		ID:       c.ID,
		Name:     c.Name,
		StartsAt: c.StartsAt,
		EndsAt:   c.EndsAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getCompetitionByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "competitionID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid competition id", http.StatusBadRequest)
		return
	}

	c, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeCompetitionError(w, err)
		return
	}

	resp := dto.CompetitionResponse{
		ID:       c.ID,
		Name:     c.Name,
		StartsAt: c.StartsAt,
		EndsAt:   c.EndsAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) joinCompetition(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "competitionID")
	competitionID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid competition id", http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.JoinCompetitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if err := h.service.JoinWithTradingAccount(r.Context(), competitionID, userID, req.Login, req.Broker, req.InvestorPassword); err != nil {
		writeCompetitionError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateAccountSize(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "competitionID")
	competitionID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid competition id", http.StatusBadRequest)
		return
	}

	accountLoginStr := chi.URLParam(r, "accountLogin")
	accountLogin, err := strconv.ParseInt(accountLoginStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateAccountSizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateAccountSize(r.Context(), competitionID, accountLogin, req.AccountSize); err != nil {
		writeCompetitionError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getLeaderboard(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "competitionID")
	competitionID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid competition id", http.StatusBadRequest)
		return
	}

	var limit int32
	var offset int32

	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = int32(n)
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
		offset = int32(n)
	}

	entries, err := h.service.GetLeaderboard(r.Context(), competitionID, limit, offset)
	if err != nil {
		writeCompetitionError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(mapper.LeaderboardToDTO(entries))
}

func (h *Handler) insertTrades(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "competitionID")
	competitionID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid competition id", http.StatusBadRequest)
		return
	}

	var req dto.InsertTradesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	trades := make([]model.Trade, 0, len(req.Trades))
	for _, it := range req.Trades {
		trades = append(trades, model.Trade{
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
		})
	}

	if err := h.service.InsertTrades(r.Context(), competitionID, req.TradingAccountLogin, trades); err != nil {
		writeCompetitionError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	competitionIDStr := chi.URLParam(r, "competitionID")
	competitionID, err := uuid.Parse(competitionIDStr)
	if err != nil {
		http.Error(w, "invalid competition id", http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	state, err := h.service.GetUserCompetitionState(r.Context(), userID, competitionID)
	if err != nil {
		http.Error(w, "failed to fetch state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(state)
}

func (h *Handler) getCurrent(w http.ResponseWriter, r *http.Request) {
	comp, err := h.service.GetCurrentCompetition(r.Context())
	if err != nil {
		http.Error(w, "competition not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(comp)
}

func (h *Handler) requestAccount(w http.ResponseWriter, r *http.Request) {
	competitionIDStr := chi.URLParam(r, "competitionID")
	competitionID, err := uuid.Parse(competitionIDStr)
	if err != nil {
		http.Error(w, "invalid competition id", http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.RequestAccount(r.Context(), userID, competitionID); err != nil {
		http.Error(w, "failed to request account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
