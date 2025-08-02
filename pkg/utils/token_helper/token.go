package token_helper

import (
	"net/http"
	"strings"
)

const (
	TokenHeader  = "Authorization"
	BearerScheme = "Bearer "

	ContextNameCompanyId    = "company_id"
	ContextNameManagedUsers = "managed_users"
	ContextNameIsBiz        = "is_biz"
	ContextNameScopes       = "scopes"
	ContextNameAccountId    = "account_id"
)

var (
	ContextHeaders = []string{
		ContextNameCompanyId,
		ContextNameManagedUsers,
		ContextNameIsBiz,
		ContextNameScopes,
		ContextNameAccountId,
	}
)

func ExtractTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get(TokenHeader)
	authHeader = strings.TrimPrefix(authHeader, BearerScheme)
	if len(authHeader) > 0 {
		return authHeader
	}
	return ""
}
