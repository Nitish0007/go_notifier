package redis

import (
	"os"
	"log"
	"time"
	"sync"
	"errors"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var once sync.Once

func InitRedis() {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
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

		_, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		} else {
			log.Println("Verified -> Connected to Redis successfully")
		}
	})
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		InitRedis()
	}
	return redisClient
}

type RedisCache struct {
	redisClient *redis.Client
	mu sync.RWMutex
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		redisClient: GetRedisClient(),
		mu: sync.RWMutex{},
	}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.redisClient.Get(ctx, key).Result()
}


func (c *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return errors.New("failed to marshal value: " + err.Error())
	}
	return c.redisClient.Set(ctx, key, jsonValue, expiration).Err()
}
