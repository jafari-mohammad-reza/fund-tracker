package services

import (
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/redis/go-redis/v9"
)

type FundService struct {
	redisClient       *redis.Client
	apiFetcherService *ApiFetcherService
}

func NewFundService() *FundService {
	redisClient := data.GetRedisClient()
	apiFetcherService := NewApiFetcher()
	return &FundService{
		redisClient:       redisClient,
		apiFetcherService: apiFetcherService,
	}
}

func (service *FundService) GetFunds() (*structs.Fund, error) {

	return nil, nil
}
