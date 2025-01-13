package handlers

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/goccy/go-json"
	"github.com/nhdms/base-go/internal"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/pkg/utils"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"github.com/spf13/viper"
	"go-micro.dev/v5/client"
	"strings"
)

type WebhookHandler struct {
	Publisher                app.PublisherInterface
	WebhookClient            services.WebhookService
	Name                     string
	webhookKindToRMQExchange map[string]string
	enableLog                bool
}

func (h *WebhookHandler) GetName() string {
	return h.Name
}

func (h *WebhookHandler) Init() error {
	h.webhookKindToRMQExchange = viper.GetStringMapString("webhook_kind_to_exchange")
	if h.webhookKindToRMQExchange == nil {
		h.webhookKindToRMQExchange = make(map[string]string)
	}
	h.enableLog = viper.GetBool("persist_log.enable")
	h.WebhookClient = internal.CreateWebhookClient(nil)
	return nil
}

func (h *WebhookHandler) HandleMessage(msg *message.Message) error {
	webhook := models.WebhookEvent{}
	_ = json.Unmarshal(msg.Payload, &webhook)

	if len(webhook.MessageUuid) == 0 {
		webhook.MessageUuid = msg.UUID
	}

	if len(webhook.MessageUuid) == 0 {
		webhook.MessageUuid = watermill.NewUUID()
	}

	isSideCarHook := len(webhook.Kind) == 0
	if isSideCarHook {
		m := make(map[string]interface{})
		_ = json.Unmarshal(msg.Payload, &m)
		webhook.Kind = internal.WebhookKindFacebookSideCar
		webhook.Payload = utils.Map2ProtoStruct(m)
	}

	// save log to db
	if isSideCarHook || h.enableLog || webhook.IsRetry {
		_, err := h.WebhookClient.InsertLogs(
			context.Background(),
			&models.WebhookEvents{Events: []*models.WebhookEvent{&webhook}},
			client.WithRetries(5),
		)

		if err != nil {
			webhook.IsRetry = true
			perr := h.Publisher.PublishRouting(h.Name, "", utils.ToJSONByte(webhook))
			if perr != nil {
				logger.DefaultLogger.Errorw("Could not push message to rmq",
					"exchange", h.Name,
					"routing_key", "",
					"error", perr.Error(),
				)
				return perr
			}
			logger.DefaultLogger.Errorw("Could not insert webhook logs", "error", err.Error(), "")
			return err
		}
	}

	logOnly := webhook.Kind == internal.WebhookKindFacebookSideCar
	if logOnly {
		return nil
	}

	exWithRoutingKey := strings.Split(h.webhookKindToRMQExchange[webhook.Kind], ":")
	if len(exWithRoutingKey) == 0 {
		return nil
	}

	exchange := exWithRoutingKey[0]
	if len(exchange) == 0 {
		return nil
	}

	routingKey := ""
	if len(exWithRoutingKey) > 1 {
		routingKey = exWithRoutingKey[1]
	}

	payload, _ := webhook.Payload.MarshalJSON()
	err := h.Publisher.PublishRouting(exchange, routingKey, payload)
	if err != nil {
		logger.DefaultLogger.Errorw("Could not push message to rmq",
			"exchange", exchange,
			"routing_key", routingKey,
			"error", err.Error(),
		)
		return err
	}

	return nil
}

func (h *WebhookHandler) SetPublisher(p app.PublisherInterface) {
	h.Publisher = p
}

func (h *WebhookHandler) Close() {
}
