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
	suite.repository.On("Count", input.Token).Return(0, nil)
	suite.repository.On("CheckLock", input.Token).Return(false, nil)
	suite.repository.On("Save", input.Token, mock.Anything).Return(nil)

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
	suite.repository.On("CheckLock", input.Ip).Return(false, nil)
	suite.repository.On("Count", input.Ip).Return(2, nil)
	suite.repository.On("LockKey", input.Ip).Return(nil)

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
	suite.repository.On("Count", input.Token).Return(2, nil)
	suite.repository.On("CheckLock", input.Token).Return(false, nil)
	suite.repository.On("LockKey", input.Token).Return(nil)

	allowed, err := suite.useCase.Execute(input)
	assert.Nil(t, err)
	assert.False(t, allowed)
}

func (suite *RateLimitTestSuite) Test_Key_Locked() {
	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "token",
		Ip:    "ip",
	}
	suite.repository.On("CheckLock", input.Token).Return(true, nil)

	allowed, err := suite.useCase.Execute(input)
	assert.Nil(t, err)
	assert.False(t, allowed)
}

func (suite *RateLimitTestSuite) Test_Error_Checking_Lock() {

	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "token",
		Ip:    "ip",
	}
	suite.repository.On("CheckLock", input.Token).Return(false, assert.AnError)

	allowed, err := suite.useCase.Execute(input)
	assert.NotNil(t, err)
	assert.False(t, allowed)
}

func (suite *RateLimitTestSuite) Test_Error_Count() {

	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "token",
		Ip:    "ip",
	}
	suite.repository.On("CheckLock", input.Token).Return(false, nil)
	suite.repository.On("Count", input.Token).Return(0, assert.AnError)

	allowed, err := suite.useCase.Execute(input)
	assert.NotNil(t, err)
	assert.False(t, allowed)
}

func (suite *RateLimitTestSuite) Test_Error_LockKey() {

	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "token",
		Ip:    "ip",
	}
	suite.repository.On("CheckLock", input.Token).Return(false, nil)
	suite.repository.On("Count", input.Token).Return(2, nil)
	suite.repository.On("LockKey", input.Token).Return(assert.AnError)

	allowed, err := suite.useCase.Execute(input)
	assert.NotNil(t, err)
	assert.False(t, allowed)
}

func (suite *RateLimitTestSuite) Test_Error_Save() {

	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "token",
		Ip:    "ip",
	}
	suite.repository.On("CheckLock", input.Token).Return(false, nil)
	suite.repository.On("Count", input.Token).Return(1, nil)
	suite.repository.On("LockKey", input.Token).Return(nil)
	suite.repository.On("Save", input.Token, mock.Anything).Return(assert.AnError)

	allowed, err := suite.useCase.Execute(input)
	assert.NotNil(t, err)
	assert.False(t, allowed)
}

func (suite *RateLimitTestSuite) Test_Token_Missing() {

	t := suite.T()
	input := usecases.RateLimitInput{
		Ip: "ip",
	}

	allowed, err := suite.useCase.Execute(input)
	assert.Nil(t, err)
	assert.False(t, allowed)
}

func (suite *RateLimitTestSuite) Test_Ip_Missing() {

	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "Token",
	}

	allowed, err := suite.useCase.Execute(input)
	assert.Nil(t, err)
	assert.False(t, allowed)
}
func TestRateLimit(t *testing.T) {
	suite.Run(t, new(RateLimitTestSuite))
}
