package handlers

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
)

type BinlogHandler struct {
	Publisher app.PublisherInterface
	Name      string
}

func (b *BinlogHandler) HandleMessage(msg *message.Message) error {
	logger.DefaultLogger.Debug("Received ", string(msg.Payload))
	return b.Publisher.PublishSimple(
		b.Name,
		msg.Payload,
	)
}

func (b *BinlogHandler) Init() error {
	return nil
}

func (b *BinlogHandler) SetPublisher(p app.PublisherInterface) {
	b.Publisher = p
}

func (b *BinlogHandler) Close() {}

func (b *BinlogHandler) GetName() string {
	return b.Name
}
