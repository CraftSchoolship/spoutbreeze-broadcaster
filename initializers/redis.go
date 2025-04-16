package initializers

import (
	"context"
	"os"
	"log"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisContext = context.Background()

func ConnectToRedis() {
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisPassword == "" {
		log.Fatalf("REDIS_PASSWORD environment variable not set")
	}

	if redisHost == "" {
		log.Fatalf("REDIS_HOST environment variable not set")
	}

	if redisPort == "" {
		log.Fatalf("REDIS_PORT environment variable not set")
	}

	redis_address := fmt.Sprintf("%s:%s",redisHost,redisPort )

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redis_address, 
		Password: redisPassword,
		DB:       0,
	})

	// Ping Redis to check connection
	_, err := RedisClient.Ping(RedisContext).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}