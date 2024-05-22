package usecases

import "context"

type RateLimitUseCase struct {
	repository string
}

type RateLimitInput struct {
	Token string
	Ip    string
}

func (u *RateLimitUseCase) Execute(input RateLimitInput) (bool, error) {
	return true, nil
}

func NewRateLimitUseCase(ctx context.Context, repository string) *RateLimitUseCase {
	return &RateLimitUseCase{repository: repository}
}
