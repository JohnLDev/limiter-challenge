package middlewares

import (
	"net"
	"net/http"

	"github.com/johnldev/rate-limiter/internal/config"
	"github.com/johnldev/rate-limiter/internal/repositories"
	usecases "github.com/johnldev/rate-limiter/internal/useCases"
)

const (
	tooManyRequestsMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func RateLimitMiddlawere(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		conf := config.GetConfig()
		useCase := usecases.NewRateLimitUseCase(r.Context(),
			// ? the use case will accept any struct that implements the Repository interface
			repositories.NewRedisRepository(r.Context(), *conf),
			*conf)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token := r.Header.Get("API_KEY")

		canGo, err := useCase.Execute(usecases.RateLimitInput{
			Token: token,
			Ip:    ip,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !canGo {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(tooManyRequestsMessage))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
