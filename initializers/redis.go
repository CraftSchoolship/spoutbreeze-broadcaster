package initializers

import (
	"context"
	"os"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisContext = context.Background()

func ConnectToRedis() {
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		log.Fatalf("REDIS_PASSWORD environment variable not set")
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Minikube Redis address
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