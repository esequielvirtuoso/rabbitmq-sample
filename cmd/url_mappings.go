package main

import (
	"fmt"

	"github.com/esequielvirtuoso/rabbitmq-sample/controllers"
	"github.com/esequielvirtuoso/rabbitmq-sample/rabbit"
	"go.uber.org/zap"
)

// mapURLs map the HTTP routes.
func mapURLs(logger *zap.Logger) {

	// TODO: manage a way to create the queue by parameter while requesting
	producer := rabbit.NewProducer(logger, queueName, connStr)
	publishController := controllers.NewPublisher(logger, producer)
	fmt.Println(publishController)

	// curl -X PUT localhost:8081/rabbit/publish
	router.PUT("/rabbit/publish", publishController.PublishMessage)
}
