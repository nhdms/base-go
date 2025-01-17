package handlers

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"gitlab.com/a7923/athena-go/internal"
    "gitlab.com/a7923/athena-go/pkg/app"
    "gitlab.com/a7923/athena-go/proto/exmsg/services"
	"gitlab.com/a7923/athena-go/pkg/logger"
)

type EventHandler struct {
    Publisher                   app.PublisherInterface
    EventClient         services.EventService
    Name                       string
}

func (h *EventHandler) GetName() string {
    return h.Name
}

func (h *EventHandler) Init() error {
	h.EventClient  = internal.CreateEventClient(nil)
    return nil
}

func (h *EventHandler) HandleMessage(msg *message.Message) error {
	logger.AthenaLogger.Debugw("Received ", "message", string(msg.Payload))
    return nil
}

func (h *EventHandler) SetPublisher(p app.PublisherInterface) {
	h.Publisher = p
}

func (h *EventHandler) Close() {}
