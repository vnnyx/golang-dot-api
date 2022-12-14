package infrastructure

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	config := NewConfig(".env")
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
	})
	return client
}
