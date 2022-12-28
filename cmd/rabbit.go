package main

import (
	"flag"
	"os"

	"github.com/esequielvirtuoso/rabbitmq-sample/handlers"
	"github.com/esequielvirtuoso/rabbitmq-sample/rabbit"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	defaultPort      = "8081"
	defaultService   = "producer-api"
	defaultQueueName = "sample"
)

var (
	connStr = getConnectionString()
	router  = gin.Default()
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("error while initializing log")
	}
	defer logger.Sync()

	var service string
	flag.StringVar(&service, "service", defaultService, "Producer API")

	if service == defaultService {
		var port string
		flag.StringVar(&port, "port", defaultPort, "Service Port")
		flag.Parse()

		mapURLs(logger)

		logger.Info("about to start the users application")
		if err := router.Run(":" + port); err != nil {
			panic(err)
		}
	}

	if service == "consumer" {
		var queueName string
		flag.StringVar(&queueName, "queue", defaultQueueName, "Queue name")
		flag.Parse()
		consume(logger, queueName)
	}

}

func consume(logger *zap.Logger, queueName string) {
	handler := handlers.New()
	consumer := rabbit.NewConsumer(logger, connStr, handler.PrintMessages)

	// Start consuming message on the specified queues
	forever := make(chan bool)

	go consumer.Consume(queueName)

	// Multiple listeners can be specified here
	<-forever
}

func getConnectionString() string {
	connStr := os.Getenv("RMQ_URL")
	if connStr == "" {
		panic("env var RMQ_URL not set")
	}
	return connStr
}
