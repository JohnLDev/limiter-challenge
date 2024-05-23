package usecases

import (
	"context"
	"log/slog"

	"github.com/johnldev/rate-limiter/internal/config"
	"github.com/johnldev/rate-limiter/internal/interfaces"
)

type RateLimitUseCase struct {
	repository interfaces.Repository
	conf       config.Config
}

type RateLimitInput struct {
	Token string
	Ip    string
}

func (u *RateLimitUseCase) Execute(input RateLimitInput) (bool, error) {
	if input.Ip == "" || input.Token == "" {
		return false, nil
	}
	isAllowed := true

	if limit, ok := u.conf.Tokens[input.Token]; ok {
		// ? verify by token
		access, err := u.repository.GetAccessByToken(input.Token)
		if err != nil {
			slog.Error(err.Error())
			return false, err
		}
		if access > limit {
			slog.Info("rate limit exceeded by token")
			isAllowed = false
		}
	} else {
		// ? verify by ip
		access, err := u.repository.GetAccessByIp(input.Ip)
		if err != nil {
			slog.Error(err.Error())
			return false, err
		}

		if access > u.conf.RateLimit {
			slog.Info("rate limit exceeded by ip")
			isAllowed = false
		}
	}
	return isAllowed, nil
}

func NewRateLimitUseCase(ctx context.Context, repository interfaces.Repository, conf config.Config) *RateLimitUseCase {
	return &RateLimitUseCase{repository: repository, conf: conf}
}
