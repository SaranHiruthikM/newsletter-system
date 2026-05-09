package api

import (
	"github.com/SaranHiruthikM/newsletter-system/internal/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, healthHandler *handlers.HealthHandler, subHandler *handlers.SubscribeHandler, confirmHandler *handlers.ConfirmHandler, newsHandler *handlers.NewsletterHandler) {
	api := app.Group("/api/v1")
	api.Get("/health", healthHandler.Check)
	api.Post("/subscribe", subHandler.Handle)
	api.Get("/confirm", confirmHandler.Handle)
	api.Post("/newsletter/send", newsHandler.Handle)
}
