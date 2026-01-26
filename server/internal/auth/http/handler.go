package http

import (
	"encoding/json"
	"github.com/filipcvejic/trading_tournament/internal/auth"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	service         *auth.AuthService
	refreshTokenTTL time.Duration
}

func NewHandler(service *auth.AuthService, refreshTokenTTL time.Duration) *Handler {
	return &Handler{service: service, refreshTokenTTL: refreshTokenTTL}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		//r.Post("/login-with-refresh", h.LoginWithRefresh)
		//r.Post("/refresh", h.Refresh)
		r.Post("/logout", h.Logout)

		// protected
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthenticationMiddleware)
			r.Get("/me", h.me)
		})
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Print(err.Error())
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	_, err := h.service.Register(r.Context(), req.Email, req.Username, req.DiscordUsername, req.Password)
	if err != nil {
		log.Print(err.Error())
		writeAuthError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	access, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
	})

	w.WriteHeader(http.StatusNoContent)
}

//func (h *Handler) LoginWithRefresh(w http.ResponseWriter, r *http.Request) {
//	var req auth.LoginRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		http.Error(w, "invalid json body", http.StatusBadRequest)
//		return
//	}
//
//	access, refresh, err := h.service.LoginWithRefresh(r.Context(), req.Email, req.Password, h.refreshTokenTTL)
//	if err != nil {
//		writeAuthError(w, err)
//		return
//	}
//
//	http.SetCookie(w, &http.Cookie{
//		Name:     "refresh_token",
//		Value:    refresh,
//		Path:     "/auth/refresh",
//		HttpOnly: true,
//		SameSite: http.SameSiteLaxMode,
//		Secure:   false,
//		MaxAge:   int(h.refreshTokenTTL.Seconds()),
//	})
//
//	w.Header().Set("Content-Type", "application/json")
//	_ = json.NewEncoder(w).Encode(auth.AuthResponse{
//		AccessToken: access,
//	})
//}
//
//func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
//	var req auth.RefreshRequest
//	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//		http.Error(w, "invalid json body", http.StatusBadRequest)
//		return
//	}
//
//	access, err := h.service.RefreshAccessToken(r.Context(), req.RefreshToken)
//	if err != nil {
//		writeAuthError(w, err)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	_ = json.NewEncoder(w).Encode(auth.AuthResponse{AccessToken: access})
//}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.Me(r.Context(), userID)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}
