package handlers

import (
    "gitlab.com/a7923/athena-go/pkg/app"
    "gitlab.com/a7923/athena-go/pkg/logger"
    transhttp "gitlab.com/a7923/athena-go/pkg/transport"
    "io"
    "net/http"
    "time"
)

type TrackingHandler struct {
    Producer app.PublisherInterface
}

func (h *TrackingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    defer func() {
        logger.AthenaLogger.Debugw("Processed request", "url", r.URL.Path, "took", time.Since(start).Milliseconds())
    }()

    requestBody, err := io.ReadAll(r.Body)
    if err != nil {
        transhttp.RespondError(w, http.StatusBadRequest, err.Error())
        return
    }

    logger.AthenaLogger.Debugw("published request", "url", r.URL.Path,"body", string(requestBody), "took", time.Since(start).Milliseconds())

    transhttp.RespondJSON(w, http.StatusOK, map[string]interface{}{
        "success": time.Now().UnixNano(),
    })
}