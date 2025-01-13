package main

import (
	app2 "github.com/nhdms/base-go/cmd/apis/webhook-api/app"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
)

func main() {
	publisher, err := app.NewPublisher()
	if err != nil {
		logger.DefaultLogger.Fatal("Start publisher failed", err)
	}
	// more connection here

	s := app2.NewServer(publisher)
	api := app.NewAPI(s)
	err = api.Run()
	if err != nil {
		logger.DefaultLogger.Fatal("Start API failed", err)
	}
}
