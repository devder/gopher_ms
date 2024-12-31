package main

import (
	"log"
	"os"

	"github.com/devder/gopher_ms/listener/event"
	"github.com/devder/gopher_ms/listener/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to rabbit-mq
	rabbitConn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to RabbitMQ")
	defer rabbitConn.Close()

	// listen for messages
	// consume
	consumer := event.NewConsumer(rabbitConn)

	// watch the queue and consume events
	consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
}
