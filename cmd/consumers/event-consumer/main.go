package main

import (
    "gitlab.com/a7923/athena-go/cmd/consumers/event-consumer/handlers"
    "gitlab.com/a7923/athena-go/pkg/app"
    "gitlab.com/a7923/athena-go/pkg/logger"
)

func main() {
    eventHandler := &handlers.EventHandler{
        Name: "event-consumer",
    }

    err := app.StartNewConsumer(eventHandler)
    if err != nil {
        logger.AthenaLogger.Fatal("Failed to start consumer: ", err)
    }
}