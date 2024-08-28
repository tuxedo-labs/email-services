package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"tuxedo-email-service/engine"
	"tuxedo-email-service/types"

	"github.com/rabbitmq/amqp091-go"
)

func Stream() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	connection, err := amqp091.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}

	ctx := context.Background()
	emailConsumer, err := channel.ConsumeWithContext(ctx, "email", "consumer-email", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}

	for message := range emailConsumer {
		var data types.EmailData

		err := json.Unmarshal(message.Body, &data)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			continue
		}
		fmt.Println("Email:", data.Email.To)
		err = engine.Send(data.Email.Subject, data.Email.Content, data.Email.From, data.Email.To)

		if err != nil {
			log.Printf("Error sending email: %v", err)
		} else {
			log.Println("Email sent successfully!")
		}
	}
}
