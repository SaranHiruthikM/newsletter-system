package queue

import "github.com/rabbitmq/amqp091-go"

func Connect(url string) (*amqp091.Connection, *amqp091.Channel, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	chann, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	return conn, chann, nil
}
