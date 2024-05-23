package usecases_test

import (
	"context"
	"testing"

	"github.com/johnldev/rate-limiter/internal/config"
	usecases "github.com/johnldev/rate-limiter/internal/useCases"
	"github.com/johnldev/rate-limiter/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RateLimitTestSuite struct {
	suite.Suite
	useCase    *usecases.RateLimitUseCase
	repository *mocks.RepositoryMock
}

func (suite *RateLimitTestSuite) SetupTest() {
	conf := config.Config{
		RateLimit: 2,
		Tokens: map[string]int{
			"token": 2,
		}}

	suite.repository = mocks.NewRepositoryMock()
	suite.useCase = usecases.NewRateLimitUseCase(context.Background(), suite.repository, conf)
}

func (suite *RateLimitTestSuite) Test_Verify_Success() {
	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "token",
		Ip:    "ip",
	}
	suite.repository.On("GetAccessByIp", mock.Anything).Return(0, nil)
	suite.repository.On("GetAccessByToken", mock.Anything).Return(0, nil)

	allowed, err := suite.useCase.Execute(input)
	assert.Nil(t, err)
	assert.True(t, allowed)
}
func (suite *RateLimitTestSuite) Test_Verify_FailByIp() {
	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "token312",
		Ip:    "ip",
	}
	suite.repository.On("GetAccessByIp", mock.Anything).Return(3, nil)
	suite.repository.On("GetAccessByToken", mock.Anything).Return(0, nil)

	allowed, err := suite.useCase.Execute(input)
	assert.Nil(t, err)
	assert.False(t, allowed)
}

func (suite *RateLimitTestSuite) Test_Verify_FailByToken() {
	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "token",
		Ip:    "ip",
	}
	suite.repository.On("GetAccessByIp", mock.Anything).Return(0, nil)
	suite.repository.On("GetAccessByToken", mock.Anything).Return(3, nil)

	allowed, err := suite.useCase.Execute(input)
	assert.Nil(t, err)
	assert.False(t, allowed)
}

func TestRateLimit(t *testing.T) {
	suite.Run(t, new(RateLimitTestSuite))
}
