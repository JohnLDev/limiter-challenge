package usecases_test

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"
	"testing"

	"github.com/johnldev/rate-limiter/internal/config"
	"github.com/johnldev/rate-limiter/internal/interfaces"
	"github.com/johnldev/rate-limiter/internal/repositories"
	usecases "github.com/johnldev/rate-limiter/internal/useCases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RateLimitIntegrationTestSuite struct {
	suite.Suite
	useCase    *usecases.RateLimitUseCase
	repository interfaces.Repository
	container  testcontainers.Container
	ctx        context.Context
}

func (suite *RateLimitIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		Name:         "redis-test",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}
	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	host, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	conf := config.Config{
		RateLimit: 100,
		DbHost:    host,
		DbPort:    portInt,
		DbName:    "0",
		BlockTime: 5000,
		Tokens: map[string]int{
			"token": 50,
		},
	}

	suite.repository = repositories.NewRedisRepository(ctx, conf)
	suite.useCase = usecases.NewRateLimitUseCase(ctx, suite.repository, conf)
	suite.container = redisC
	suite.ctx = ctx
}

func (s *RateLimitIntegrationTestSuite) TearDownSuite() {
	if err := s.container.Terminate(s.ctx); err != nil {
		log.Fatalf("Could not stop redis: %s", err)
	}
}

func (s *RateLimitIntegrationTestSuite) Test101RequestsIp() {
	t := s.T()
	const requestNumber = 101
	input := usecases.RateLimitInput{
		Token: "token2",
		Ip:    "ip",
	}
	wg := new(sync.WaitGroup)
	wg.Add(requestNumber)
	var success = 0
	var failed = 0
	for i := 0; i < requestNumber; i++ {
		go func(i int) {
			defer wg.Done()
			isAllowed, err := s.useCase.Execute(input)
			assert.Nil(t, err)
			if isAllowed {
				success++
			} else {
				failed++
			}
		}(i)
	}
	wg.Wait()
	assert.Equal(t, 100, success)
	assert.Equal(t, 1, failed)
}

func (s *RateLimitIntegrationTestSuite) Test500RequestsToken() {
	t := s.T()
	const requestNumber = 500
	input := usecases.RateLimitInput{
		Token: "token",
		Ip:    "ip",
	}
	wg := new(sync.WaitGroup)
	wg.Add(requestNumber)
	var success = 0
	var failed = 0
	for range requestNumber {
		go func() {
			defer wg.Done()
			isAllowed, err := s.useCase.Execute(input)
			assert.Nil(t, err)
			if isAllowed {
				success++
			} else {
				failed++
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, 50, success)
	assert.Equal(t, 450, failed)
}

func TestRateLimitIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimitIntegrationTestSuite))
}
