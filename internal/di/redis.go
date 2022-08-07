package di

import (
	"jadwalin-auth-events-integrator/internal/shared/config"

	"github.com/go-redis/redis/v8"
)

func NewRedis(config *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       config.Redis.Database,
	})
}
