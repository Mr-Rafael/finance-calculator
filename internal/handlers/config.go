package handlers

import (
	"net/http"
	"sync/atomic"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
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
