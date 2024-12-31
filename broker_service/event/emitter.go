package event

import (
	"log"

	"github.com/devder/gopher_ms/broker/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.conn.Channel()
	util.FailOnError(err, "Failed to open channel")
	defer channel.Close()

	return declareExchange(channel)
}

func (e *Emitter) Push(event, severity string) error {
	ch, err := e.conn.Channel()
	util.FailOnError(err, "Failed to open channel")
	defer ch.Close()

	log.Println("[x] Pushing event", event, "with severity", severity)

	return ch.Publish(
		"logs_topic", // exchange
		severity,     // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		})
}

func NewEventEmitter(conn *amqp.Connection) Emitter {
	emitter := Emitter{
		conn: conn,
	}

	err := emitter.setup()
	util.FailOnError(err, "Failed to declare exchange")

	return emitter
}
