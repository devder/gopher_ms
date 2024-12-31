package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to rabbit-mq
	rabbitConn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to RabbitMQ")
	defer rabbitConn.Close()

	// listen for messages

	// consume

	// watch the queue and consume events
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
