package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devder/gopher_ms/listener/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Reuse the HTTP client
var client = &http.Client{
	Timeout: 10 * time.Second,
}

type Consumer struct {
	conn *amqp.Connection
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) Consumer {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	util.FailOnError(err, "Failed to declare exchange")

	return consumer
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	util.FailOnError(err, "Failed to open channel")

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) {
	ch, err := consumer.conn.Channel()
	util.FailOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	util.FailOnError(err, "Failed to declare random queue")

	// bind the queue to the exchange
	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"logs_topic", // exchange
			false,        // no-wait
			nil,          // args
		)

		util.FailOnError(err, "Failed to bind queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var payload Payload
			err := json.Unmarshal(d.Body, &payload)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}
			go handlePayload(payload)
		}
	}()

	log.Printf(" [*] Waiting for messages [Exchange, Queue]. [logs_topic, %s]\n", q.Name)
	<-forever // block forever
}

func handlePayload(payload Payload) {
	log.Printf("Handling payload: Name=%s, Data=%s", payload.Name, payload.Data)
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		logEvent(payload)
	case "auth":
		// do something
	default:
		// handle default
	}
}

func logEvent(entry Payload) {
	// send json to microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// call the service
	req, err := http.NewRequest(http.MethodPost, "http://logger/log", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		return
	}

	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling logger service: %v", fmt.Errorf("%w", err))
		return
	}

	defer resp.Body.Close()

	// ensure we get back the correct status
	if resp.StatusCode != http.StatusAccepted {
		errMsg := "failed to call logger service"
		log.Printf("%s: %v", errMsg, err)
	}
}
