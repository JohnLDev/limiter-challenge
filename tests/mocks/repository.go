package mocks

import "github.com/stretchr/testify/mock"

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAccessByToken(token string) (int, error) {
	args := r.Called(token)
	return args.Int(0), args.Error(1)
}

func (r *RepositoryMock) GetAccessByIp(ip string) (int, error) {
	args := r.Called(ip)
	return args.Int(0), args.Error(1)
}

func NewRepositoryMock() *RepositoryMock {
	return new(RepositoryMock)
}
