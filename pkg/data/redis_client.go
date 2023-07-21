package data

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

var redisClient *redis.Client

func SetupRedisClient() error {

	redisHost, _ := strconv.Atoi(os.Getenv("REDIS_HOST"))
	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redisDialTimeout, _ := strconv.Atoi(os.Getenv("REDIS_DIAL_TIMEOUT"))
	redisReadTimeout, _ := strconv.Atoi(os.Getenv("REDIS_READ_TIMEOUT"))
	redisWriteTimeout, _ := strconv.Atoi(os.Getenv("REDIS_WRITE_TIMEOUT"))
	redisPoolSize, _ := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
	redisPoolTimeout, _ := strconv.Atoi(os.Getenv("REDIS_POOL_TIMEOUT"))

	redisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:           0,
		DialTimeout:  time.Duration(redisDialTimeout) * time.Second,
		ReadTimeout:  time.Duration(redisReadTimeout) * time.Second,
		WriteTimeout: time.Duration(redisWriteTimeout) * time.Second,
		PoolSize:     redisPoolSize,
		PoolTimeout:  time.Duration(redisPoolTimeout) * time.Second,
	})
	return nil
}
func CloseRedisClient() {
	redisClient.Close()
}
func GetRedisClient() *redis.Client {
	return redisClient
}
