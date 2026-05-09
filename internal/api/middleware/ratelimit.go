package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client, limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		ctx := context.Background()
		key := fmt.Sprintf("rate_limit:%s", ip)
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{"error": "something went wrong"},
			)
		}

		if count == 1 {
			rdb.Expire(ctx, key, window)
		} else if count > int64(limit) {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				fiber.Map{"error": "too many requests"},
			)
		}

		return c.Next()
	}
}
