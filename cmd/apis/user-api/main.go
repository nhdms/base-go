package main

import (
	app2 "github.com/nhdms/base-go/cmd/apis/user-api/app"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
)

func main() {
	s := &app2.Server{}
	api := app.NewAPI(s)
	err := api.Run()
	if err != nil {
		logger.DefaultLogger.Fatal("Start API failed", err)
	}
}
