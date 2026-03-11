package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/golang-jwt/jwt/v5"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Queries        *db.Queries
	AccessSecret   string
	RefreshSecret  string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			respondWithErrorCode(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(cfg.AccessSecret), nil
		})

		if err != nil || !token.Valid {
			respondWithErrorCode(w, "expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
