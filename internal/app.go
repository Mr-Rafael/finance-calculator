package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mr-Rafael/finance-calculator/internal/api"
	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/repository"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type App struct {
	Handler http.Handler
	DB      *pgxpool.Pool
}

func New() *App {
	ctx := context.Background()

	// environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error reading .env: %v", err)
		return &App{}
	}
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// database
	dbURL := os.Getenv("POSTGRES_CONNECTION_STRING")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(pool)

	// repos
	usersRepo := repository.NewUsersRepo(queries)
	authRepo := repository.NewAuthRepo(queries)

	// services
	userService := service.NewUserService(&usersRepo)
	authService := service.NewAuthService(&authRepo, &usersRepo, accessSecret, refreshSecret)

	// handlers
	adminHandler := api.NewAdminHandler()
	userHandler := api.NewUsersHandler(userService)
	authHandler := api.NewAuthHandler(authService)

	// middlewares

	// mux
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", adminHandler.HandlerHealthZ)
	mux.HandleFunc("POST /app/users/create", userHandler.CreateUser)
	mux.HandleFunc("POST /app/login", authHandler.Login)

	return &App{
		Handler: mux,
		DB:      pool,
	}
}

func (a *App) Run() {
	defer a.DB.Close()
	port := ":8080"
	http.ListenAndServe(port, a.Handler)
}
