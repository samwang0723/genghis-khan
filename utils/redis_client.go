package utils

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	linkType := os.Getenv("REDIS_TYPE")
	redisURL := os.Getenv("REDIS_URL")
	password := ""
	if !strings.Contains(linkType, "docker") && !strings.Contains(redisURL, "localhost") {
		parsedURL, _ := url.Parse(redisURL)
		password, _ = parsedURL.User.Password()
		redisURL = parsedURL.Host
	}
	log.Printf("Connecting to %s", redisURL)
	client = redis.NewClient(&redis.Options{
		Addr:         redisURL,
		Password:     password,
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
