package transhttp

import (
	"net/http"
	"strings"
)

func IsWebSocket(r *http.Request) bool {
	// Check if the "Connection" header contains "Upgrade" (case-insensitive)
	if !strings.EqualFold(r.Header.Get("Connection"), "Upgrade") {
		return false
	}

	// Check if the "Upgrade" header contains "websocket" (case-insensitive)
	if !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return false
	}

	return true
}
