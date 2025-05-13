package config

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	cacheInstance *redis.Client
	cacheOnce     sync.Once
)

// GetCache Returns Redis instance
func GetCache(ctx context.Context) *redis.Client {
	cacheOnce.Do(func() {
		host := GetEnv("REDIS_HOST", "localhost")
		port := GetEnv("REDIS_PORT", "6379")
		password := GetEnv("REDIS_PASSWORD", "")
		db := 0

		cacheInstance = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       db,
		})
	})
	if err := cacheInstance.Ping(ctx).Err(); err != nil {
		cacheInstance = nil
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	return cacheInstance
}

// CloseCache Closes Redis connection
func CloseCache() error {
	if cacheInstance == nil {
		return nil
	}
	return cacheInstance.Close()
}
