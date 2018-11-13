package utils

import (
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_URL"),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	if _, err := client.Ping().Result(); err != nil {
		log.Fatal("Cannot open redis connection")
	}
}

func RedisClient() *redis.Client {
	return client
}
