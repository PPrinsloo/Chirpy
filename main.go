package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

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

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))

	handler := Chain(fileServer, Logger, AddHeader)
	mux.Handle("/", handler)

	fmt.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
