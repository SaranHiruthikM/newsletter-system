package main

import (
	"github.com/SaranHiruthikM/newsletter-system/internal/api"
	"github.com/SaranHiruthikM/newsletter-system/internal/api/handlers"
	"github.com/SaranHiruthikM/newsletter-system/internal/config"
	"github.com/SaranHiruthikM/newsletter-system/internal/database"
	"github.com/SaranHiruthikM/newsletter-system/internal/queue"
	"github.com/SaranHiruthikM/newsletter-system/internal/redis"
	"github.com/SaranHiruthikM/newsletter-system/internal/repository/postgres"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DB)
	if err != nil {
		panic("DB dbection error")
	}

	subRepo := postgres.NewSubscriberRepository(db)
	newsRepo := postgres.NewNewsletterRepository(db)

	_, channel, err := queue.Connect(cfg.RabbitMQ.URL)
	if err != nil {
		panic("failed to connect to RabbitMQ: " + err.Error())
	}
	publisher := queue.NewPublisher(channel)

	healthHandler := handlers.NewHealthHandler(db)
	subHandler := handlers.NewSubscriberHandler(subRepo, publisher)
	confirmHandler := handlers.NewConfirmHandler(subRepo)
	newsHandler := handlers.NewNewsletterHandler(subRepo, newsRepo, publisher)

	redisClient := redis.Connect(cfg.Redis)
	adminKey := cfg.App.AdminKey

	app := fiber.New()
	api.SetupRoutes(app, healthHandler, subHandler, confirmHandler, newsHandler, redisClient, cfg, adminKey)

	app.Listen(":" + cfg.App.Port)

}
