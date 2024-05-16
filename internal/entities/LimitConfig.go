package entities

import "errors"

const (
	IP    = "ip"
	Token = "token"
)

type LimitConfig struct {
	Key  string
	Max  int
	Type string
}

func NewLimitConfig(key string, max int, t string) (*LimitConfig, error) {
	if t != IP && t != Token {
		return nil, errors.New("invalid type")
	}

	return &LimitConfig{
		Key:  key,
		Max:  max,
		Type: t,
	}, nil
}
