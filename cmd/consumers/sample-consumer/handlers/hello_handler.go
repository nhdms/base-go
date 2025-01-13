package handlers

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"strings"
	"time"
)

type HelloHandler struct {
	UserClient services.UserService
	Publisher  app.PublisherInterface
	Name       string
}

func (h *HelloHandler) GetName() string {
	return h.Name
}

func (h *HelloHandler) Init() error {
	return nil
}

func (h *HelloHandler) HandleMessage(msg *message.Message) error {
	logger.DefaultLogger.Infof("Received message: %s", msg.Payload)
	start := time.Now()
	if strings.Contains(string(msg.Payload), "test") {
		logger.DefaultLogger.Infof("start sleep for 10 seconds...  (workerId: %v)  ", msg.UUID)
		time.Sleep(10 * time.Second)
	}
	logger.DefaultLogger.Infof("done took %v", time.Since(start))
	return nil
}

func (h *HelloHandler) SetPublisher(p app.PublisherInterface) {
	h.Publisher = p
}

func (h *HelloHandler) Close() {
}
