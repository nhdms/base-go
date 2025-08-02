package app

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/goccy/go-json"
	config2 "github.com/nhdms/base-go/pkg/config"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/pkg/replication/pglogicalstream"
	"os"
	"os/signal"
	"syscall"
)

func StartDataReplica(handler Consumer) error {
	name := handler.GetName()
	streamConfig := &pglogicalstream.Config{}
	err := config2.LoadConfigToVar(streamConfig, "postgres")
	if err != nil {
		return err
	}

	err = handler.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize gRPC client for task %s: %w", name, err)
	}

	defer handler.Close()
	// Capture OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	publisher, err := NewPublisher()
	if err != nil {
		return fmt.Errorf("failed to create publisher for task %s: %w", name, err)
	}
	defer publisher.Close()

	//publisher.Publish("", message.NewMessage(name, []byte("Hello, Athena!")))
	// init redis....
	handler.SetPublisher(publisher)

	pgStream, err := pglogicalstream.NewPgStream(streamConfig)
	if err != nil {
		panic(err)
	}

	// Listen for shutdown signal
	go func() {
		<-signalChan
		logger.DefaultLogger.Infof("Shutdown signal received")
		_ = pgStream.Stop()
		_ = publisher.Close()
		handler.Close()
	}()

	pgStream.OnMessage(func(changeCaptured pglogicalstream.Wal2JsonChanges) {
		msgBytes, _ := json.Marshal(changeCaptured)
		newMsg := message.NewMessage(watermill.NewUUID(), msgBytes)
		err = handler.HandleMessage(newMsg)
		if err != nil {
			logger.DefaultLogger.Error("Failed to handle message", err, watermill.LogFields{
				"uuid": newMsg.UUID,
			})
			// stop to prevent message being lost
			err = pgStream.Stop()
			if err != nil {
				logger.DefaultLogger.Fatal("Failed to stop streaming", err)
			}
		}
		if changeCaptured.Lsn != nil {
			// Snapshots dont have LSN
			pgStream.AckLSN(*changeCaptured.Lsn)
		}
	})
	return nil
}
