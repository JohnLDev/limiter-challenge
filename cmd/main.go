package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/johnldev/rate-limiter/internal/config"
)

func main() {
	slog.Info("Starting server", slog.String("service_name", config.Conf.ServiceName))

	mux := http.NewServeMux()

	mux.HandleFunc("POST /{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.PathValue("id"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe(":8080", mux)
}
