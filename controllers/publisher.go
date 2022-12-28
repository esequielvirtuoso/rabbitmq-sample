package controllers

import (
	"fmt"
	"net/http"

	"github.com/esequielvirtuoso/rabbitmq-sample/rabbit"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
)

type Producer interface {
	PublishMessage(contentType, queue string, body []byte)
}

type Log interface {
	Error(string, ...zapcore.Field)
}

type Publisher struct {
	log      Log
	producer Producer
}

func (p *Publisher) PublishMessage(c *gin.Context) {
	var msg rabbit.Message

	request_id := c.GetString("x-request-id")
	queue := c.GetHeader("queue")

	// bind request payload with the message model
	if binderr := c.ShouldBindJSON(&msg); binderr != nil {
		p.log.Error(fmt.Sprintf("%s, %s", binderr.Error(), request_id))

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	p.producer.PublishMessage("text/plain", queue, []byte(msg.Message))

	c.JSON(http.StatusOK, gin.H{
		"response": "Message received",
	})
}

func NewPublisher(log Log, producer Producer) *Publisher {
	return &Publisher{log: log,
		producer: producer,
	}
}
