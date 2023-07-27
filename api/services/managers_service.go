package services

import (
	"github.com/jafari-mohammad-reza/fund-tracker/pkg/data"
	"github.com/redis/go-redis/v9"
)

type ManagersService struct {
	redisClient       *redis.Client
	apiFetcherService *ApiFetcherService
}

func NewManagersService() *ManagersService {
	redisClient := data.GetRedisClient()
	apiFetcherService := NewApiFetcher()
	return &ManagersService{
		redisClient,
		apiFetcherService,
	}
}
