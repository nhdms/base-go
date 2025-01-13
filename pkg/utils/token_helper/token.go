package token_helper

import (
	"net/http"
	"strings"
)

const (
	TokenHeader  = "Authorization"
	BearerScheme = "Bearer "
)

func ExtractTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get(TokenHeader)
	authHeader = strings.TrimPrefix(authHeader, BearerScheme)
	if len(authHeader) > 0 {
		return authHeader
	}
	return ""
}
