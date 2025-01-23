package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conn *amqp091.Connection
var Channel *amqp091.Channel
var DBCollection *mongo.Collection

// connection and initialise
func Connect() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading environment variables", err)
	}

	rabbit_uri := os.Getenv("RABBITMQ_URI")
    mongoURI := os.Getenv("MONGO_URI")
    mongoDBName := os.Getenv("MONGO_DB")
    mongoCollectionName := os.Getenv("MONGO_COLLECTION")

	var err error
	conn, err := amqp091.Dial(rabbit_uri)
	if err != nil {
		return fmt.Errorf("failed connection to RabbitMQ: %v", err)
	}

	Channel, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create Channel: %v", err)
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        return fmt.Errorf("failed to connect to MongoDB: %v", err)
    }

	DBCollection = client.Database(mongoDBName).Collection(mongoCollectionName)
    
	return nil
}

func PublishToQueue(queueName string, jobID interface{}) error {
	if Channel == nil {
		return fmt.Errorf("channel is not initialized")
	}

	_, err := Channel.QueueDeclare(
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
	
	// jobID to string if needed
	body := fmt.Sprintf("%v", jobID)
	err = Channel.Publish(
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
