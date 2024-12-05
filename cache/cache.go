package cache

import (
	"context"
	"time"
	"x-straight-check/config"
	"x-straight-check/log"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Env().Cache.Address,
		Password: config.Env().Cache.Password,
		DB:       0,
	})

}

func SetWithTTL(key, value string, ttl time.Duration) error {
	err := redisClient.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Errorf("setting the key \"%s\" with the value \"%s\" failed.", key, value)
		return err
	}
	return nil
}

func Get(key string) (value string, err error) {
	value, err = redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Errorf("Key \"%s\" was not found. Key is either expired or not registered.", key)
		return
	} else if err != nil {
		return
	}
	return
}

func Del(key string) error {
	_, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		log.Errorf("Key \"%s\" was not deleted. Key is either expired or not registered.", key)
		return err
	}
	return nil
}
