package redisservice

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *redis.Client
	Ctx    context.Context
}

var redisService *RedisService

func InitRedisService(redisAddr, redisPassword string, db int) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       db,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	redisService = &RedisService{
		Client: client,
		Ctx:    ctx,
	}
}

func ReadFromCache(key string) (string, error) {
	val, err := redisService.Client.Get(redisService.Ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key '%s' does not exist in cache", key)
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func WriteToCache(key, value string) error {
	err := redisService.Client.Set(redisService.Ctx, key, value, time.Minute).Err()
	if err != nil {
		return err
	}
	log.Printf("Key '%s' written to cache with expiration: 1 MIN", key)
	return nil
}

func DeleteCache(key string) error {
	err := redisService.Client.Del(redisService.Ctx, key).Err()
	if err != nil {
		return err
	}
	log.Printf("Key '%s' deleted from cache", key)
	return nil
}
