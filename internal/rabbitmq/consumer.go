package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"

	"mini-avito/internal/avito_service"
)

type StatusConsumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	repo avito_service.AdRepository
}

type adStatusMessage struct {
	AdID   string `json:"ad_id"`
	Status string `json:"status"`
}

func NewStatusConsumer(amqpURL string, repo avito_service.AdRepository) (*StatusConsumer, error) {
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
		"ad_status_changed",
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

	return &StatusConsumer{
		conn: conn,
		ch:   ch,
		repo: repo,
	}, nil
}

func (c *StatusConsumer) Start(ctx context.Context) {
	msgs, err := c.ch.Consume(
		"ad_status_changed",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[Consumer] Failed to register consumer: %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var incoming adStatusMessage
				if err := json.Unmarshal(msg.Body, &incoming); err != nil {
					log.Printf("[Consumer] Error unmarshaling JSON: %v\n", err)
					msg.Nack(false, false)
					continue
				}

				adUUID, err := uuid.Parse(incoming.AdID)
				if err != nil {
					log.Printf("[Consumer] Error parsing UUID: %v\n", err)
					msg.Nack(false, false)
					continue
				}

				if err := c.repo.UpdateAdStatus(ctx, adUUID, incoming.Status); err != nil {
					log.Printf("[Consumer] Error updating ad status in DB: %v\n", err)
					msg.Nack(false, true)
					continue
				}

				log.Printf("[Consumer] Successfully updated ad %s status to '%s'\n", incoming.AdID, incoming.Status)
				msg.Ack(false)
			}
		}
	}()
}

func (c *StatusConsumer) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	log.Println("[Consumer] Connection closed")
}
