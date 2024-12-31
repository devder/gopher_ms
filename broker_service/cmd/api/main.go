package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/devder/gopher_ms/broker/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	rabbitConn *amqp.Connection
}

func main() {
	// connect to rabbit-mq
	rabbitConn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to RabbitMQ")
	defer rabbitConn.Close()

	app := Config{
		rabbitConn: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
