package token

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nhdms/base-go/proto/exmsg/models"
)

type Token struct {
	*jwt.Token
}

type Processor interface {
	GetToken(ctx context.Context, tokenString string) (token *Token, err error)
	ExtractMetadata(token *Token) map[string]string
	CheckPermissions(token *Token, requirePermissions map[int64]int64) bool
	GenerateToken(ctx context.Context, user *models.User) (string, error)
}
