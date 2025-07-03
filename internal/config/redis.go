package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// NewRedisConnection creates a singleton Redis connection
func NewRedisConnection(cfg *Config) (*redis.Client, error) {
	if redisClient != nil {
		return redisClient, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	log.Println("Redis connection is perfect!")
	redisClient = client
	return redisClient, nil
}

// GetRedisClient returns the singleton Redis client
func GetRedisClient() *redis.Client {
	return redisClient
}
