package middleware

import "github.com/gofiber/fiber/v2"

func APIKeyAuth(adminKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")

		if apiKey == "" || apiKey != adminKey {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{"error": "unauthorized access"},
			)
		}

		return c.Next()
	}
}
