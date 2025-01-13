package handlers

import (
	"github.com/goccy/go-json"
	"github.com/nhdms/base-go/internal"
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
	transhttp "github.com/nhdms/base-go/pkg/transport"
	"github.com/nhdms/base-go/pkg/utils"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"io"
	"net/http"
	"time"
)

type FacebookWebhookHandler struct {
	Producer app.PublisherInterface
}

func (h *FacebookWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		logger.DefaultLogger.Debugw("Processed request", "url", r.URL.Path, "took", time.Since(start).Milliseconds())
	}()
	err := h.verifySignature(r)
	if err != nil {
		transhttp.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		transhttp.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	m := make(map[string]interface{})
	_ = json.Unmarshal(requestBody, &m)

	message := models.WebhookEvent{
		Kind:    internal.WebhookKindFacebook,
		Payload: utils.Map2ProtoStruct(m),
		Meta:    nil, // header here if
	}

	startPublished := time.Now()
	err = h.Producer.PublishRoutingPersist(
		internal.WebhookExchange,
		"",
		utils.ToJSONByte(message),
	)
	logger.DefaultLogger.Debugw("published request", "url", r.URL.Path, "took", time.Since(startPublished).Milliseconds())

	if err != nil {
		transhttp.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	transhttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": time.Now().UnixNano(),
	})
}

func (h *FacebookWebhookHandler) verifySignature(r *http.Request) error {
	return nil
}
