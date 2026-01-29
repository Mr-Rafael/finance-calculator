package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mr-Rafael/finance-calculator/internal/handlers"
)

func main() {
	port := ":8080"
	mux := http.NewServeMux()

	var config handlers.ApiConfig
	config.FileserverHits.Store(0)

	mux.Handle("/app/", config.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./files")))))
	mux.HandleFunc("GET /api/healthz", handlers.HandlerHealthZ)
	mux.HandleFunc("GET /admin/metrics", config.HandlerMetrics)
	mux.HandleFunc("GET /app/savings/calculate", config.HandlerCalculateGet)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	fmt.Printf("Starting server on %v\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
