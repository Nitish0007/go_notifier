package redis_utils

import (
	"context"
	"log"
	"time"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() *redis.Client {
	if RedisClient != nil {
		return RedisClient
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
		PoolSize: 10,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MinIdleConns: 5,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Println("Connected to Redis")
			return nil
		},
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return RedisClient
}

func GetRedisClient() *redis.Client {
	if RedisClient == nil {
		return ConnectRedis()
	}
	return RedisClient
}

func SetRedis(ctx context.Context, key string, value any, expiration time.Duration) error {
	client := GetRedisClient()
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Error setting value in Redis: %v", err)
		return err
	}
	return nil
}

func GetRedis(ctx context.Context, key string) (string, error) {
	client := GetRedisClient()
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Error getting value from Redis: %v", err)
		return "", err
	}
	return value, nil
}

func DeleteRedis(ctx context.Context, key string) error {
	client := GetRedisClient()
	err := client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error deleting value from Redis: %v", err)
		return err
	}
	return nil
}

func SetRedisJSON(ctx context.Context, key string, value any, expiration time.Duration) error {
	client := GetRedisClient()
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Error setting value in Redis: %v", err)
		return err
	}
	return nil
}

func GetRedisJSON(ctx context.Context, key string) (any, error) {
	client := GetRedisClient()
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Error getting value from Redis: %v", err)
		return nil, err
	}
	return value, nil
}	

func DeleteRedisJSON(ctx context.Context, key string) error {
	client := GetRedisClient()
	err := client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error deleting value from Redis: %v", err)
		return err
	}
	return nil
}

// func SetRedisHash(ctx context.Context, key string, field string, value any, expiration time.Duration) error {
// 	client := GetRedisClient()
// 	err := client.HSet(ctx, key, field, value).Err()
// 	if err != nil {
// 		log.Printf("Error setting value in Redis: %v", err)
// 		return err
// 	}
// 	return nil
// }

// func GetRedisHash(ctx context.Context, key string, field string) (any, error) {
// 	client := GetRedisClient()
// 	value, err := client.HGet(ctx, key, field).Result()
// 	if err != nil {
// 		log.Printf("Error getting value from Redis: %v", err)
// 		return "", err
// 	}
// 	return value, nil
// }