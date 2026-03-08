package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	port := ":8080"
	mux := http.NewServeMux()
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error reading .env: %v", err)
		return
	}

	var config handlers.ApiConfig
	config.FileserverHits.Store(0)
	config.AccessSecret = os.Getenv("ACCESS_SECRET")
	config.RefreshSecret = os.Getenv("REFRESH_SECRET")

	ctx := context.Background()

	dbURL := os.Getenv("POSTGRES_CONNECTION_STRING")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	config.Queries = db.New(pool)

	mux.Handle("/app/", config.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./files")))))
	mux.HandleFunc("GET /api/healthz", handlers.HandlerHealthZ)
	mux.HandleFunc("GET /admin/metrics", config.HandlerMetrics)
	mux.HandleFunc("POST /app/savings/calculate", config.HandlerSavingsCalculateGet)
	mux.HandleFunc("POST /app/loans/calculate", config.HandlerLoansCalculateGet)
	mux.HandleFunc("POST /app/users/create", config.HandlerUsersCreate)
	mux.HandleFunc("POST /app/users/login", config.HandlerUsersLogin)
	mux.HandleFunc("POST /app/users/refresh", config.HandlerRefresh)

	server := &http.Server{
		Addr:    port,
		Handler: withCORS(mux),
	}

	fmt.Printf("Starting server on %v\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func withCORS(h http.Handler) http.Handler {
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}
