package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("access_token")
		if err != nil {
			log.Println("AUTH: missing access_token cookie:", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		secret := []byte(os.Getenv("JWT_SECRET"))
		if len(secret) == 0 {
			log.Println("AUTH: JWT_SECRET is empty")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// ✅ enforce HS256/HMAC family
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %T", t.Method)
			}
			return secret, nil
		})

		if err != nil {
			log.Println("AUTH: jwt parse error:", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Println("AUTH: token not valid")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("AUTH: claims type assert failed")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			log.Println("AUTH: missing/invalid sub claim:", claims["sub"])
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			log.Println("AUTH: sub is not uuid:", sub, "err:", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID retrieves the user ID from the request context
func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

//func AuthenticationMiddleware(authService *AuthService) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			authHeader := r.Header.Get("Authorization")
//			if authHeader == "" {
//				http.Error(w, "Authorization header required", http.StatusUnauthorized)
//				return
//			}
//
//			parts := strings.Split(authHeader, " ")
//			if len(parts) != 2 || parts[0] != "Bearer" {
//				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
//				return
//			}
//
//			tokenString := parts[1]
//
//			claims, err := authService.ValidateToken(tokenString)
//			if err != nil {
//				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
//				return
//			}
//
//			userIDStr, ok := claims["sub"].(string)
//			if !ok {
//				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
//				return
//			}
//
//			userID, err := uuid.Parse(userIDStr)
//			if err != nil {
//				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
//				return
//			}
//
//			ctx := context.WithValue(r.Context(), UserIDKey, userID)
//
//			next.ServeHTTP(w, r.WithContext(ctx))
//		})
//	}
//}
