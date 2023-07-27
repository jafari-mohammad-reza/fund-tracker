package data

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
	"time"
)

var redisClient *redis.Client

func SetupRedisClient() error {

	//redisHost := os.Getenv("REDIS_HOST")
	//redisPort := os.Getenv("REDIS_PORT")
	redisDialTimeout, _ := strconv.Atoi(os.Getenv("REDIS_DIAL_TIMEOUT"))
	redisReadTimeout, _ := strconv.Atoi(os.Getenv("REDIS_READ_TIMEOUT"))
	redisWriteTimeout, _ := strconv.Atoi(os.Getenv("REDIS_WRITE_TIMEOUT"))
	redisPoolSize, _ := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
	redisPoolTimeout, _ := strconv.Atoi(os.Getenv("REDIS_POOL_TIMEOUT"))

	redisClient = redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6380",
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
func SetValue(ctx context.Context, client *redis.Client, key string, value []byte, duration time.Duration) error {
	client.Set(ctx, key, value, duration)
	return nil
}
func GetValue(ctx context.Context, client *redis.Client, key string) (string, error) {
	cachedData, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return cachedData, nil
}

func GetDataFromCacheOrFetch[T any](fetchFunc func() (*T, error), key string, ctx context.Context, redisClient *redis.Client) (*T, error) {
	dataStr, err := GetValue(ctx, redisClient, key)
	if err != nil {
		if err != redis.Nil {
			// Some error occurred while fetching from Redis, but not because the key is not found
			return nil, err
		}
		// The key is not found in Redis, fetch the data using the provided fetch function
		responseData, err := fetchFunc()
		if err != nil {
			return nil, err
		}
		// Set the fetched data in Redis
		dataJSON, err := json.Marshal(responseData)
		if err != nil {
			log.Printf("Error marshalling data: %v\n", err)
		} else {
			err := SetValue(ctx, redisClient, key, dataJSON, time.Hour*3)
			if err != nil {
				return nil, err
			}
		}
		return responseData, nil
	} else {
		// Key found in Redis, unmarshal the data
		var responseData T
		err := json.Unmarshal([]byte(dataStr), &responseData)
		if err != nil {
			log.Printf("Error unmarshalling data: %v\n", err)
			return nil, err
		}
		return &responseData, nil
	}
}
