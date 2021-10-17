package repository

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

func NewRedisClient(config RedisConfig) *redis.Client {
	redisHost := config.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: config.MinIdleConns,
		PoolSize:     config.PoolSize,
		PoolTimeout:  time.Duration(config.PoolTimeout) * time.Second,
		Password:     config.Password,
		DB:           config.DB,
	})

	return client
}
