package http

import (
	"encoding/json"
	"github.com/filipcvejic/trading_tournament/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	service *user.Service
}

func NewHandler(service *user.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.createUser)
		r.Get("/{userID}", h.getUserByID)
	})
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	u, err := h.service.Create(r.Context(), req.Email, req.Username, req.DiscordUsername, req.Password)
	if err != nil {
		writeUserError(w, err)
		return
	}

	resp := user.UserResponse{
		ID:        u.ID.String(),
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	u, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeUserError(w, err)
		return
	}

	resp := user.UserResponse{
		ID:              u.ID.String(),
		Email:           u.Email,
		Username:        u.Username,
		DiscordUsername: u.DiscordUsername,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
