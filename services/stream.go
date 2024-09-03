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
	connection, channel := setupRabbitMQ(rabbitMQURL)
	defer connection.Close()
	defer channel.Close()

	ctx := context.Background()
	consumeEmails(ctx, channel)
}

func setupRabbitMQ(rabbitMQURL string) (*amqp091.Connection, *amqp091.Channel) {
	connection, err := amqp091.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Error creating channel: %v", err)
	}

	return connection, channel
}

func consumeEmails(ctx context.Context, channel *amqp091.Channel) {
	emailConsumer, err := channel.ConsumeWithContext(ctx, "email", "consumer-email", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}

	for message := range emailConsumer {
		processEmailMessage(message)
	}
}

func processEmailMessage(message amqp091.Delivery) {
	var data types.EmailData

	if err := json.Unmarshal(message.Body, &data); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return
	}

	fmt.Println("Sending email to:", data.Email.To)
	if err := engine.Send(data.Email.Subject, data.Email.Content, data.Email.From, data.Email.To); err != nil {
		log.Printf("Error sending email: %v", err)
	} else {
		log.Println("Email sent successfully!")
	}
}
