package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

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
		true,
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
			log.Printf("[Worker] Received a message: %s", string(msg.Body))
		}
	}()

	log.Println("[Worker] Started listening for messages. To exit press CTRL+C")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Worker] Shutting down...")
}