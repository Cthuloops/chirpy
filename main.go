package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/",
		apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/",
		http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("GET /healthz", readiness)

	mux.HandleFunc("GET /metrics", apiCfg.serveMetrics)
	mux.HandleFunc("POST /reset", apiCfg.reset)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func readiness(writer http.ResponseWriter, reader *http.Request) {
	_ = reader
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) serveMetrics(writer http.ResponseWriter,
	reader *http.Request) {
	_ = reader
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	writer.Write([]byte(hits))
}

func (cfg *apiConfig) reset(writer http.ResponseWriter,
	reader *http.Request) {
	_ = reader
	writer.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}
