package main

import (
	"github.com/nhdms/base-go/cmd/services/webhook-service/handlers"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/dbtool"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/proto/exmsg/services"
)

func main() {
	svc := app.NewGRPCService()
	psql, err := dbtool.NewConnectionManager(dbtool.DBTypePostgreSQL, nil)
	if err != nil {
		logger.DefaultLogger.Fatal("Failed to connect to database: ", err)
	}

	grpcSvc := handlers.NewWebhookHandler(psql)
	err = services.RegisterWebhookServiceHandler(svc.Server(), grpcSvc)
	if err != nil {
		logger.DefaultLogger.Fatal(err)
	}

	err = svc.Run()
	if err != nil {
		logger.DefaultLogger.Fatal(err)
	}
}
