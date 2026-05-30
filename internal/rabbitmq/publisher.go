package rabbitmq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewPublisher(amqpURL string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"ad_created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	log.Println("[RabbitMQ] Successfully connected and declared 'ad_created' queue")

	return &Publisher{
		conn: conn,
		ch:   ch,
	}, nil
}

func (p *Publisher) PublishAdCreated(ctx context.Context, adID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.ch.PublishWithContext(ctx,
		"",
		"ad_created",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(adID),
		})

	if err != nil {
		log.Printf("[RabbitMQ] Error publishing message: %v\n", err)
		return err
	}

	log.Printf("[RabbitMQ] ID %s successfully published to queue\n", adID)
	return nil
}

func (p *Publisher) Close() {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	log.Println("[RabbitMQ] Connection closed")
}
