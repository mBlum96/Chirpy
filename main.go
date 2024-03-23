package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	appRouter := chi.NewRouter()
	apiRouter := chi.NewRouter()
	apiCfg := &apiConfig{
		fileserverHits: 0,
	}
	apiRouter.Mount("/api")
	appsMiddleWare := middlewareCors(apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	appRouter.Handle("/app/*", middlewareCors(appsMiddleWare))
	appRouter.Handle("/app", middlewareCors(appsMiddleWare))
	appRouter.Get("/metrics", apiCfg.handleMetrics)
	appRouter.Handle("/reset", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits = 0
		w.WriteHeader(http.StatusOK)
	}))
	appRouter.Get("/healthz", handleHealthz)
	corsMux := middlewareCors(appRouter)
	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}
	server.ListenAndServe()
}
