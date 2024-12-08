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

func InitRedisService(redisAddr string) *RedisService {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return &RedisService{
		Client: client,
		Ctx:    ctx,
	}
}

func (redisService *RedisService) ReadFromCache(key string) (string, error) {
	val, err := redisService.Client.Get(redisService.Ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key '%s' does not exist in cache", key)
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (redisService *RedisService) WriteToCache(key, value string, durationUntilNextMinute time.Duration) error {
	now := time.Now()
	nextFullMinute := now.Truncate(time.Minute).Add(time.Minute)
	if durationUntilNextMinute == 0 {
		durationUntilNextMinute = nextFullMinute.Sub(now)
	}
	err := redisService.Client.Set(redisService.Ctx, key, value, durationUntilNextMinute).Err()
	if err != nil {
		return err
	}
	log.Printf("Key '%s' written to cache with expiration: 1 MIN", key)
	return nil
}

func (redisService *RedisService) DeleteCache(key string) error {
	err := redisService.Client.Del(redisService.Ctx, key).Err()
	if err != nil {
		return err
	}
	log.Printf("Key '%s' deleted from cache", key)
	return nil
}
