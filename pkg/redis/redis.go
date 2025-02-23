package redis

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
	once   sync.Once
)

func InitRedis(addr, password string) {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		})

		ctx := context.Background()
		_, err := client.Ping(ctx).Result()
		if err != nil {
			log.Fatal("Error connecting to Redis:", err)
		}
	})
}

func GetClient() *redis.Client {
	if client == nil {
		log.Fatal("Redis client not initialized")
	}
	return client
}
