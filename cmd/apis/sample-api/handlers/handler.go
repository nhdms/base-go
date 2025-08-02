package handlers

import (
	"github.com/nhdms/base-go/pkg/app"
	"github.com/nhdms/base-go/pkg/logger"
	transhttp "github.com/nhdms/base-go/pkg/transport"
	"io"
	"net/http"
	"time"
)

type SampleHandler struct {
	Producer app.PublisherInterface
}

func (h *SampleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		logger.DefaultLogger.Debugw("Processed request", "url", r.URL.Path, "took", time.Since(start).Milliseconds())
	}()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		transhttp.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.DefaultLogger.Debugw("published request", "url", r.URL.Path, "body", string(requestBody), "took", time.Since(start).Milliseconds())

	transhttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": time.Now().UnixNano(),
	})
}
