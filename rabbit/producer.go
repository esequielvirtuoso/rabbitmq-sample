package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.uber.org/zap/zapcore"
)

type ProducerLog interface {
	Error(string, ...zapcore.Field)
}

type Producer struct {
	connectionString string
	log              ProducerLog
}

func (p *Producer) onError(err error, queue, msg string) {
	if err != nil {
		p.log.Error(fmt.Sprintf("Error occurred while publishing message on '%s' queue. Error message: %s", queue, msg))
	}
}

func (p *Producer) PublishMessage(contentType, queue string, body []byte) {
	conn, err := amqp.Dial(p.connectionString)
	p.onError(err, queue, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	p.onError(err, queue, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	p.onError(err, queue, "failed to declare queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
	p.onError(err, queue, "failed to publish a message")
}

func NewProducer(log ProducerLog, connStr string) *Producer {
	return &Producer{log: log,
		connectionString: connStr}
}
