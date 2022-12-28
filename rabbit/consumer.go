package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.uber.org/zap/zapcore"
)

type ConsummerLog interface {
	Error(string, ...zapcore.Field)
	Info(string, ...zapcore.Field)
}

type Consumer struct {
	connectionString string
	log              ConsummerLog
	handleMessage    func(queue string, msg amqp.Delivery, err error)
}

func (c *Consumer) onError(err error, queue, msg string) {
	if err != nil {
		c.handleMessage(queue, amqp.Delivery{}, err)
	}
}

func (c *Consumer) Consume(queue string) {
	conn, err := amqp.Dial(c.connectionString)
	c.onError(err, queue, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	c.onError(err, queue, "failed to open channel")

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	c.onError(err, queue, "failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	c.onError(err, queue, "failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			c.handleMessage(queue, d, nil)
		}
	}()

	c.log.Info(fmt.Sprintf("Started listening for messages on '%s' queue", queue))
	<-forever
}

func NewConsumer(log ConsummerLog, connStr string, handler func(queue string, msg amqp.Delivery, err error)) *Consumer {
	return &Consumer{log: log,
		connectionString: connStr,
		handleMessage:    handler}
}
