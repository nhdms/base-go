package main

import (
	"github.com/nhdms/base-go/cmd/consumers/data-replica/binlog-consumer/handlers"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
)

func main() {
	handler := &handlers.BinlogHandler{
		Name: "binlog_consumer",
	}

	err := app.StartDataReplica(handler)
	if err != nil {
		logger.DefaultLogger.Fatal("Failed to start data replica consumer: ", err)
	}
}
