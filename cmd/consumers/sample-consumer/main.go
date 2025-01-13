package main

import (
	"github.com/nhdms/base-go/cmd/consumers/sample-consumer/handlers"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
)

func main() {
	helloHandler := &handlers.HelloHandler{
		Name: "hello_iam_go",
	}

	err := app.StartNewConsumer(helloHandler)
	if err != nil {
		logger.DefaultLogger.Fatal("Failed to start consumer: ", err)
	}
}
