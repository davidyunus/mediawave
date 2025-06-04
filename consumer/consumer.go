package consumer

import (
	"encoding/json"
	"log"

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

func processOrder(order Order) error {
	log.Printf("Processing order ID: %s", order.OrderID)
	// Simulasi proses
	log.Println("Payment processed")
	log.Println("Inventory updated")
	log.Println("Email sent to", order.Email)
	log.Println("Shipping label generated")
	return nil
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

	err = ch.Qos(1, 0, false) // satu pesan sekali proses (prefetch=1)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	msgs, err := ch.Consume(
		"order_queue",
		"",    // consumer tag
		false, // autoAck: false = manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var order Order
			if err := json.Unmarshal(d.Body, &order); err != nil {
				log.Println("⚠ Failed to decode message:", err)
				d.Nack(false, false) // drop message, masuk DLQ kalau ada
				continue
			}

			if err := processOrder(order); err != nil {
				log.Printf("⚠ Error processing order %s: %v", order.OrderID, err)
				d.Nack(false, true) // requeue
				continue
			}

			d.Ack(false)
		}
	}()

	log.Println("Loading...")
	<-forever
}
