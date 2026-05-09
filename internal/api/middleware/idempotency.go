package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func Idempotency(rdb *redis.Client, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idemKey := c.Get("Idempotency-Key")
		if idemKey == "" {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{"error": "idempotency key is missing"},
			)
		}

		lockKey := fmt.Sprintf("idempotency:lock:%s", idemKey)
		responseKey := fmt.Sprintf("idempotency:response:%s", idemKey)
		ctx := context.Background()

		cachedResponse, err := rdb.Get(ctx, responseKey).Result()
		if err != nil && err != redis.Nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				fiber.Map{"error": "something went wrong"},
			)
		}

		if err == nil {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.Status(fiber.StatusOK).Send([]byte(cachedResponse))
		}

		isLocked, _ := rdb.SetNX(ctx, lockKey, "locked", ttl).Result()
		if !isLocked {
			return c.Status(fiber.StatusTooManyRequests).JSON(
				fiber.Map{"error": "request already in progress"},
			)
		}

		if err := c.Next(); err != nil {
			return err
		}

		rdb.Set(ctx, responseKey, c.Response().Body(), ttl)
		return nil
	}
}
