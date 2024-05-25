package mocks

import (
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Save(key string, id string) error {
	args := r.Called(key, id)
	return args.Error(0)
}

func (r *RepositoryMock) Count(key string) (int, error) {
	args := r.Called(key)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) CheckLock(key string) (bool, error) {
	args := r.Called(key)
	return args.Bool(0), args.Error(1)
}

func (r *RepositoryMock) LockKey(key string) error {
	args := r.Called(key)
	return args.Error(0)
}

func NewRepositoryMock() *RepositoryMock {
	return new(RepositoryMock)
}
