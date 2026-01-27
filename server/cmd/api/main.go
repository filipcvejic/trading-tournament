package main

import (
	"github.com/filipcvejic/trading_tournament/db"
	"github.com/filipcvejic/trading_tournament/internal/auth"
	authhttp "github.com/filipcvejic/trading_tournament/internal/auth/http"
	"github.com/filipcvejic/trading_tournament/internal/competition"
	competitionhttp "github.com/filipcvejic/trading_tournament/internal/competition/http"
	"github.com/filipcvejic/trading_tournament/internal/tradingaccount"
	tradingaccounthttp "github.com/filipcvejic/trading_tournament/internal/tradingaccount/http"
	"github.com/filipcvejic/trading_tournament/internal/user"
	userhttp "github.com/filipcvejic/trading_tournament/internal/user/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	requiredVars := []string{"DATABASE_URL"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Required environment variable %s is not set", v)
		}
	}
}

func main() {
	loadEnv()

	database := db.NewDatabase(os.Getenv("DATABASE_URL"))

	competitionRepo := competition.NewPostgresRepository(database)
	competitionService, err := competition.NewService(competitionRepo, os.Getenv("CRYPTO_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	competitionHandler := competitionhttp.NewHandler(competitionService)

	userRepo := user.NewPostgresRepository(database)
	userService := user.NewService(userRepo)
	userHandler := userhttp.NewHandler(userService)

	tradingAccountRepo := tradingaccount.NewPostgresRepository(database)
	tradingAccountService := tradingaccount.NewService(tradingAccountRepo)
	tradingAccountHandler := tradingaccounthttp.NewHandler(tradingAccountService)

	refreshTokenRepo := auth.NewPostgresRefreshTokenRepository(database)
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, os.Getenv("JWT_SECRET"), 15)
	authHandler := authhttp.NewHandler(authService, 60)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://main.d11w7ewevo758h.amplifyapp.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300, // 5 min
	}))

	competitionHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)
	tradingAccountHandler.RegisterRoutes(r)
	authHandler.RegisterRoutes(r)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
