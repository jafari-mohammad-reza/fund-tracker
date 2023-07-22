package services

import (
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/structs"
	"github.com/redis/go-redis/v9"
)

type FundService struct {
	redisClient *redis.Client
}

func NewFundService() *FundService {
	redisClient := data.GetRedisClient()
	return &FundService{
		redisClient: redisClient,
	}
}

func (service *FundService) GetFunds() (*structs.Fund, error) {
	return nil, nil
}
