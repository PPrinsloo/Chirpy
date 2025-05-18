package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

// Middleware defines a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// Chain applies multiple middleware to a handler in sequence
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	// Apply middleware in reverse order (from last to first)
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed: %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func AddHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Chirpy")
		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) checkMetrics(w http.ResponseWriter, _ *http.Request) {
	hitCount := cfg.fileserverHits.Load()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", hitCount)
}

func (cfg *apiConfig) reset(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Swap(0)
	w.WriteHeader(http.StatusOK)
}

func routs(mux *http.ServeMux, cfg *apiConfig) {

	fileServer := http.FileServer(http.Dir("."))
	// For fileServer, we still need to use Handle since it's an http.Handler
	mux.Handle("/app/", Chain(fileServer, Logger, AddHeader, cfg.metricsInc))

	// Convert the other routes to use HandleFunc with middleware
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		handler := Chain(http.HandlerFunc(healthCheckHandler), Logger, AddHeader)
		handler.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		handler := Chain(http.HandlerFunc(cfg.checkMetrics), Logger, AddHeader)
		handler.ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /reset", func(w http.ResponseWriter, r *http.Request) {
		handler := Chain(http.HandlerFunc(cfg.reset), Logger, AddHeader)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	cfg := &apiConfig{}

	routs(mux, cfg)

	fmt.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
