package usecases_test

import (
	"context"
	"testing"

	usecases "github.com/johnldev/rate-limiter/internal/useCases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RateLimitTestSuite struct {
	suite.Suite
	useCase *usecases.RateLimitUseCase
}

func (suite *RateLimitTestSuite) SetupTest() {
	suite.useCase = usecases.NewRateLimitUseCase(context.Background(), "repository")
}

func (suite *RateLimitTestSuite) Test_Verify_Success() {
	t := suite.T()
	input := usecases.RateLimitInput{
		Token: "token",
		Ip:    "ip",
	}
	_, err := suite.useCase.Execute(input)
	assert.Nil(t, err)
}
func TestRateLimit(t *testing.T) {
	suite.Run(t, new(RateLimitTestSuite))
}
