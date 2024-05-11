package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type PurchaseNotification struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Product string `json:"product"`
	Amount  int    `json:"amount"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
		"pubsub",
		false, // durable
		false, //delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msg, err := ch.Consume(
		q.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to consume a message: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msg {
			log.Printf("Received a message: %s", d.Body)
			var purchase PurchaseNotification
			if err := json.Unmarshal(d.Body, &purchase); err != nil {
				log.Printf("Failed to parse JSON message: %v", err)
				continue
			}
			notification := map[string]interface{}{
				"data": map[string]interface{}{
					"id":      purchase.ID,
					"user":    purchase.User,
					"product": purchase.Product,
					"amount":  purchase.Amount,
				},
				"message": "Purchase success",
			}
			notificationJSON, err := json.Marshal(notification)
			if err != nil {
				log.Printf("Failed to marshal notification JSON: %v", err)
				continue
			}
			log.Printf("Formatted notification: %s", notificationJSON)
		}

	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
