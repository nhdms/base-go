package main

import (
    "gitlab.com/a7923/athena-go/cmd/services/event-service/handlers"
    "gitlab.com/a7923/athena-go/pkg/app"
    "gitlab.com/a7923/athena-go/pkg/dbtool"
    "gitlab.com/a7923/athena-go/pkg/logger"
    "gitlab.com/a7923/athena-go/proto/exmsg/services"
)

func main() {
    svc := app.NewGRPCService()
    psql, err := dbtool.NewConnectionManager(dbtool.DBTypePostgreSQL, nil)
    if err != nil {
        logger.AthenaLogger.Fatal("Failed to connect to database: ", err)
    }

    grpcSvc := handlers.NewEventHandler(psql)
    err = services.RegisterEventServiceHandler(svc.Server(), grpcSvc)
    if err != nil {
        logger.AthenaLogger.Fatal(err)
    }

    err = svc.Run()
    if err != nil {
        logger.AthenaLogger.Fatal(err)
    }
}