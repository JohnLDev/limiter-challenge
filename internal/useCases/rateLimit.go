package usecases

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/docker/distribution/uuid"
	"github.com/johnldev/rate-limiter/internal/config"
	"github.com/johnldev/rate-limiter/internal/interfaces"
)

var mutex = new(sync.Mutex)

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
	mutex.Lock()
	defer mutex.Unlock()

	isAllowed := true
	var key = input.Ip
	var limit = u.conf.RateLimit

	if tokenLimit, ok := u.conf.Tokens[input.Token]; ok {
		// ? verify by token
		limit = tokenLimit
		key = input.Token
	}

	blocked, err := u.repository.CheckLock(key)
	if err != nil {
		slog.Error(err.Error())
		return false, err
	}
	if blocked {
		slog.Info(fmt.Sprintf("key %s is blocked", key))
		return false, nil
	}

	access, err := u.repository.Count(key)
	if err != nil {
		slog.Error(err.Error())
		return false, err
	}

	if access >= limit {
		slog.Info("rate limit exceeded")
		err = u.repository.LockKey(key)

		if err != nil {
			slog.Error(err.Error())
			return false, err
		}

		isAllowed = false
	}

	if isAllowed {
		err = u.repository.Save(key, uuid.Generate().String())
		if err != nil {
			slog.Error(err.Error())
			return false, err
		}
	}

	return isAllowed, nil
}

func NewRateLimitUseCase(ctx context.Context, repository interfaces.Repository, conf config.Config) *RateLimitUseCase {
	return &RateLimitUseCase{repository: repository, conf: conf}
}
