package handlers

import (
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/nhdms/base-go/internal"
	tracking_provider "github.com/nhdms/base-go/internal/tracking-provider"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/dbtool"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/pkg/utils"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"github.com/spf13/cast"
)

type EventHandler struct {
	Publisher       app.PublisherInterface
	EventClient     services.EventService
	Name            string
	MarketingClient services.MarketingService
	allProviders    []tracking_provider.EventProvider
	RedisClient     *redis.Client
	AgsaleClient    services.AGSaleService
}

func (h *EventHandler) GetName() string {
	return h.Name
}

func (h *EventHandler) Init() error {
	var err error
	h.RedisClient, err = dbtool.CreateRedisConnection(nil)
	if err != nil {
		return err
	}

	h.EventClient = internal.CreateEventClient(nil)
	h.MarketingClient = internal.CreateMarketingClient(nil)
	h.AgsaleClient = internal.CreateAGSaleClient(nil)

	h.allProviders = []tracking_provider.EventProvider{
		&tracking_provider.LandingTracking{
			MarketingClient: h.MarketingClient,
			AgsaleClient:    h.AgsaleClient,
			RedisClient:     h.RedisClient,
		},
	}

	return nil
}

func (h *EventHandler) HandleMessage(msg *message.Message) error {
	logger.DefaultLogger.Debugw("Received ", "message", string(msg.Payload))
	event := &models.Event{}
	_ = json.Unmarshal(msg.Payload, event)
	if len(event.Type) == 0 {
		return nil
	}

	listProviders := h.filterProviders(event)
	if len(listProviders) == 0 {
		return nil
	}

	retryProviders := make([]string, 0)

	for _, p := range listProviders {
		err := p.SendTracking(msg.Context(), event)
		if err != nil {
			retryProviders = append(retryProviders, p.GetName())
		}
	}

	if len(retryProviders) > 0 {
		meta := event.Metadata.AsMap()
		if meta == nil {
			meta = make(map[string]interface{})
		}
		meta["retry_providers"] = strings.Join(retryProviders, ",")
		event.Metadata = utils.Map2ProtoStruct(meta)
		err := h.Publisher.PublishSimple(
			h.Name,
			utils.ToJSONByte(event))

		if err != nil {
			logger.DefaultLogger.Errorw("Could not publish event",
				"event_type", event.Type,
				"error", err.Error())
			return err
		}
		return nil
	}

	return nil
}

func (h *EventHandler) SetPublisher(p app.PublisherInterface) {
	h.Publisher = p
}

func (h *EventHandler) Close() {}

func (h *EventHandler) filterProviders(event *models.Event) []tracking_provider.EventProvider {
	providers := make([]tracking_provider.EventProvider, 0)
	retries := make(map[string]bool)
	eventMetadata := event.GetMetadata().AsMap()
	for _, providerName := range strings.Split(cast.ToString(eventMetadata["retry_providers"]), ",") {
		if len(providerName) == 0 {
			continue
		}
		retries[providerName] = true
	}

	for _, p := range h.allProviders {
		if len(retries) > 0 {
			if !retries[p.GetName()] {
				continue
			}
			providers = append(providers, p)
			continue
		}

		if utils.StringSliceContains(p.ListEventTracking(), event.Type) {
			providers = append(providers, p)
		}
	}

	return providers
}
