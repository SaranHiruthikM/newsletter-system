package redis

import (
	"github.com/SaranHiruthikM/newsletter-system/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Host + ":" + cfg.Port,
	})
}
