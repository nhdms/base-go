package token

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/nhdms/base-go/proto/exmsg/services"
)

type SimpleTokenProcessor struct {
}

func NewSimpleTokenProcessor(*redis.Client, services.UserService) *SimpleTokenProcessor {
	return &SimpleTokenProcessor{}
}

func (s SimpleTokenProcessor) GetToken(ctx context.Context, tokenString string) (token *Token, err error) {
	//TODO implement me
	panic("implement me")
}

func (s SimpleTokenProcessor) ExtractMetadata(token *Token) map[string]string {
	//TODO implement me
	panic("implement me")
}

func (s SimpleTokenProcessor) CheckPermissions(token *Token, requirePermissions map[int64]int64) bool {
	//TODO implement me
	panic("implement me")
}

func (s SimpleTokenProcessor) GenerateToken(ctx context.Context, user *models.User) (string, error) {
	//TODO implement me
	panic("implement me")
}
