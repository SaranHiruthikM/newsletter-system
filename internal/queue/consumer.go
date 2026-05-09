package queue

import (
	"encoding/json"
	"log"
	"time"

	"github.com/SaranHiruthikM/newsletter-system/internal/email"
	"github.com/SaranHiruthikM/newsletter-system/internal/metrics"
	"github.com/SaranHiruthikM/newsletter-system/internal/repository"
	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	channel        *amqp091.Channel
	emailProvider  email.Provider
	newsletterRepo repository.NewsletterRepository
}

func NewConsumer(channel *amqp091.Channel, emailProvider email.Provider, newsletterRepo repository.NewsletterRepository) *Consumer {
	return &Consumer{
		channel:        channel,
		emailProvider:  emailProvider,
		newsletterRepo: newsletterRepo,
	}
}

func (c *Consumer) ConsumeConfirmations() error {
	_, err := c.channel.QueueDeclare("confirmation.queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgCh, err := c.channel.Consume("confirmation.queue", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for val := range msgCh {
			var payload ConfirmationPayload
			if err := json.Unmarshal(val.Body, &payload); err != nil {
				val.Nack(false, true)
				continue
			}

			log.Printf("processing confirmation for: %s", payload.Email)

			start := time.Now()
			err := c.emailProvider.Send(
				payload.Email,
				"Confirm your subscription",
				"Click here to confirm: http://localhost:8080/api/v1/confirm?token="+payload.Token,
			)
			duration := time.Since(start).Seconds()
			metrics.EmailProcessingDuration.Observe(duration)

			if err != nil {
				metrics.EmailsFailed.Inc()
				log.Printf("failed to send email to %s: %v", payload.Email, err)
				val.Nack(false, true)
			} else {
				metrics.EmailsSent.Inc()
				val.Ack(false)
			}
		}
	}()

	return nil

}

func (c *Consumer) ConsumeNewsletters() error {
	_, err := c.channel.QueueDeclare("newsletter.queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgCh, err := c.channel.Consume("newsletter.queue", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for val := range msgCh {
			var payload NewsletterPayload
			if err := json.Unmarshal(val.Body, &payload); err != nil {
				val.Nack(false, true)
				continue
			}

			start := time.Now()
			err := c.emailProvider.Send(
				payload.Email,
				payload.Subject,
				payload.Body,
			)
			duration := time.Since(start).Seconds()
			metrics.EmailProcessingDuration.Observe(duration)

			if err != nil {
				metrics.EmailsFailed.Inc()
				log.Printf("failed to sendx to %s: %v", payload.Email, err)
				val.Nack(false, true) // requeue
				c.newsletterRepo.IncrementFailCount(payload.NewsletterID)
			} else {
				metrics.EmailsSent.Inc()
				val.Ack(false)
				c.newsletterRepo.IncrementSentCount(payload.NewsletterID)
			}
		}
	}()

	return nil
}
