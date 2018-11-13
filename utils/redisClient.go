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
	herokuURL := os.Getenv("REDIS_URL")
	password := ""
	if !strings.Contains(herokuURL, "localhost") {
		parsedURL, _ := url.Parse(herokuURL)
		password, _ = parsedURL.User.Password()
		herokuURL = parsedURL.Host
	}
	log.Printf("Connecting to %s", herokuURL)
	client = redis.NewClient(&redis.Options{
		Addr:         herokuURL,
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
