package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Order struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Items   []Item `json:"items"`
}

type Item struct {
	SKU string `json:"sku"`
	Qty int    `json:"qty"`
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"order_queue",
		true,  // durable
		false, // delete when unused
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	order := Order{
		OrderID: "ORD456",
		UserID:  "USR123",
		Email:   "user@example.com",
		Items:   []Item{{SKU: "ABC123", Qty: 2}},
	}

	body, _ := json.Marshal(order)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Println("Order sent:", string(body))
}
