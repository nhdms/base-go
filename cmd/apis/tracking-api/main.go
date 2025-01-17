package main

import (
    app2 "gitlab.com/a7923/athena-go/cmd/apis/tracking-api/app"
    "gitlab.com/a7923/athena-go/pkg/app"
    "gitlab.com/a7923/athena-go/pkg/logger"
)

func main() {
    publisher, err := app.NewPublisher()
    if err != nil {
        logger.AthenaLogger.Fatal("Start publisher failed", err)
    }
    defer publisher.Close()

    s := app2.NewServer(publisher)
    api := app.NewAPI(s)
    err = api.Run()
    if err != nil {
        logger.AthenaLogger.Fatal("Start API failed", err)
    }
}