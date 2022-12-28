package handlers

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.uber.org/zap/zapcore"
)

type HandlerLog interface {
	Error(string, ...zapcore.Field)
	Info(string, ...zapcore.Field)
}

type Handler struct {
	log HandlerLog
}

func (h *Handler) PrintMessages(queue string, msg amqp.Delivery, err error) {
	if err != nil {
		h.log.Error(fmt.Sprintf("error occurred in RMQ consumer due to %s", err.Error()))
	}

	h.log.Info(fmt.Sprintf("Message received on '%s' queue: %s", queue, string(msg.Body)))
}

func New() *Handler {
	return &Handler{}
}
