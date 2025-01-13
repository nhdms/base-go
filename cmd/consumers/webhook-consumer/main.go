package main

import (
	"github.com/nhdms/base-go/cmd/consumers/webhook-consumer/handlers"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
)

func main() {
	webhookHandler := &handlers.WebhookHandler{
		Name: "webhook_consumer",
	}

	err := app.StartNewConsumer(webhookHandler)
	if err != nil {
		logger.DefaultLogger.Fatal("Failed to start consumer: ", err)
	}
}
