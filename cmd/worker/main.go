package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type StatusUpdateMsg struct {
	AdID   string `json:"ad_id"`
	Status string `json:"status"`
}

func main() {
	amqpURL := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("[Worker] Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[Worker] Failed to open a channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"ad_created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[Worker] Failed to declare 'ad_created' queue: %v", err)
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
		log.Fatalf("[Worker] Failed to declare 'ad_status_changed' queue: %v", err)
	}

	msgs, err := ch.Consume(
		"ad_created",
		"",
		false, 
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[Worker] Failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			adID := string(msg.Body)
			log.Printf("[Worker] Processing ad ID: %s", adID)

			time.Sleep(5 * time.Second)

			statuses := []string{"approved", "rejected"}
			newStatus := statuses[rand.Intn(len(statuses))]

			updateMsg := StatusUpdateMsg{
				AdID:   adID,
				Status: newStatus,
			}

			body, err := json.Marshal(updateMsg)
			if err != nil {
				log.Printf("[Worker] JSON marshal error: %v", err)
				msg.Nack(false, false)
				continue
			}

			err = ch.PublishWithContext(context.Background(),
				"",
				"ad_status_changed",
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        body,
				})

			if err != nil {
				log.Printf("[Worker] Failed to publish status update: %v", err)
				msg.Nack(false, true)
				continue
			}

			log.Printf("[Worker] Successfully processed and published status '%s' for ad %s", newStatus, adID)
			
			msg.Ack(false)
		}
	}()

	log.Println("[Worker] Started listening for messages. To exit press CTRL+C")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Worker] Shutting down...")
}