package main

import (
	"github.com/nhdms/base-go/cmd/consumers/event-consumer/handlers"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
)

func main() {
	eventHandler := &handlers.EventHandler{
		Name: "event-consumer",
	}

	err := app.StartNewConsumer(eventHandler)
	if err != nil {
		logger.DefaultLogger.Fatal("Failed to start consumer: ", err)
	}
}
