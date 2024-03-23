package main

import (
	"fmt"
	"net/http"
	"os"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("metrics.html")
	if err != nil {
		http.Error(w, "Error reading message format", http.StatusInternalServerError)
		return
	}
	metricsMessage := fmt.Sprintf(string(data), cfg.fileserverHits)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metricsMessage))
}
