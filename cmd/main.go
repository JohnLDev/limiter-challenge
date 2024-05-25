package main

import (
	"net/http"

	"github.com/johnldev/rate-limiter/internal/middlewares"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /", middlewares.RateLimitMiddlawere(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})))

	http.ListenAndServe(":8080", mux)
}
