package services

import (
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/redis/go-redis/v9"
)

type FundInfoService struct {
	redisClient *redis.Client
}

func NewFundInfoService() *FundInfoService {
	redisClient := data.GetRedisClient()
	return &FundInfoService{redisClient}
}
