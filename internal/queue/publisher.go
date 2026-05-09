package queue

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	chann *amqp091.Channel
}

func NewPublisher(chann *amqp091.Channel) *Publisher {
	return &Publisher{
		chann: chann,
	}
}

func (p *Publisher) PublishConfirmation(payload ConfirmationPayload) error {
	_, err := p.chann.QueueDeclare("confirmation.queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.chann.PublishWithContext(
		context.Background(),
		"",
		"confirmation.queue",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

}

func (p *Publisher) PublishNewsletter(payload NewsletterPayload) error {
	_, err := p.chann.QueueDeclare("newsletter.queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.chann.PublishWithContext(
		context.Background(),
		"",
		"newsletter.queue",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
