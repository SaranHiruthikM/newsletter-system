package api

import (
	"github.com/SaranHiruthikM/newsletter-system/internal/api/handlers"
	"github.com/SaranHiruthikM/newsletter-system/internal/api/middleware"
	"github.com/SaranHiruthikM/newsletter-system/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(app *fiber.App, healthHandler *handlers.HealthHandler, subHandler *handlers.SubscribeHandler, confirmHandler *handlers.ConfirmHandler, newsHandler *handlers.NewsletterHandler, redisClient *redis.Client, cfg *config.Config, adminKey string) {
	api := app.Group("/api/v1")
	api.Get("/health", healthHandler.Check)

	if cfg.RateLimit.Enabled {
		public := api.Group("", middleware.RateLimiter(redisClient, cfg.RateLimit.Limit, cfg.RateLimit.Window))
		public.Post("/subscribe", subHandler.Handle)
		public.Get("/confirm", confirmHandler.Handle)
	} else {
		api.Post("/subscribe", subHandler.Handle)
		api.Get("/confirm", confirmHandler.Handle)
	}

	var protected fiber.Router

	if cfg.Idempotency.Enabled {
		protected = api.Group("", middleware.APIKeyAuth(adminKey), middleware.Idempotency(redisClient, cfg.Idempotency.TTL))
	} else {
		protected = api.Group("", middleware.APIKeyAuth(adminKey))
	}
	protected.Post("/newsletter/send", newsHandler.Handle)
}
