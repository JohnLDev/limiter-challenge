package main

import (
	"fmt"
	"net/http"

	"github.com/johnldev/rate-limiter/internal/config"
)

func main() {

	mux := http.NewServeMux()

	fmt.Println(config.GetConfig())
	mux.HandleFunc("POST /{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.PathValue("id"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe(":8080", mux)
}
