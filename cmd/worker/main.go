package main

import (
	"log"
	"net/http"

	"github.com/SaranHiruthikM/newsletter-system/internal/config"
	"github.com/SaranHiruthikM/newsletter-system/internal/database"
	"github.com/SaranHiruthikM/newsletter-system/internal/email"
	"github.com/SaranHiruthikM/newsletter-system/internal/queue"
	"github.com/SaranHiruthikM/newsletter-system/internal/repository/postgres"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DB)
	if err != nil {
		panic("DB dbection error")
	}

	newsRepo := postgres.NewNewsletterRepository(db)

	_, channel, err := queue.Connect(cfg.RabbitMQ.URL)
	if err != nil {
		panic("failed to connect to RabbitMQ: " + err.Error())
	}

	provider := email.NewResendProvider(cfg.Email)

	consumer := queue.NewConsumer(channel, provider, newsRepo)
	consumer.ConsumeConfirmations()
	consumer.ConsumeNewsletters()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":"+cfg.App.WorkerMetricsPort, nil)
	}()

	forever := make(chan struct{})
	log.Println("worker is running, waiting for messages...")
	<-forever
}
