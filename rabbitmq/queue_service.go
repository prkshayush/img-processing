package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection
var channel *amqp091.Channel

// connection and initialise
func Connect() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading environment variables", err)
	}

	rabbit_uri := os.Getenv("RABBITMQ_URI")

	var err error
	conn, err := amqp091.Dial(rabbit_uri)
	if err != nil {
		return fmt.Errorf("failed connection to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create Channel: %v", err)
	}

	return nil
}

func PublishToQueue(queueName string, jobID interface{}) error {
	if channel == nil {
		return fmt.Errorf("channel is not initialized")
	}

	_, err := channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Auto delete
		false,     // Exclusive
		false,     // No wait
		nil,       // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}
	
	// Convert jobID to string if needed
	body := fmt.Sprintf("%v", jobID)
	err = channel.Publish(
		"",           // Default exchange
		queueName,    // Queue name
		false,        // Mandatory
		false,        // Immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	return nil
}
