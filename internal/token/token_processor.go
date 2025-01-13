package token

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nhdms/base-go/internal/permissions"
	"github.com/nhdms/base-go/pkg/common"
	"github.com/nhdms/base-go/pkg/config"
	"github.com/nhdms/base-go/pkg/utils"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/nhdms/base-go/proto/exmsg/services"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
	"time"
)

type TokenProcessor struct {
	rd                   *redis.Client
	userService          services.UserService
	jwtSecret            []byte
	enableSignatureCheck bool
}

type Claim struct {
	UserInfo []byte `json:"u"`
	jwt.MapClaims
}

func NewTokenProcessor(rd *redis.Client, userService services.UserService) *TokenProcessor {
	return &TokenProcessor{
		rd:                   rd,
		userService:          userService,
		jwtSecret:            []byte(viper.GetString("jwt.secret")),
		enableSignatureCheck: viper.GetBool("jwt.enable_signature_check"),
	}
}

func (p *TokenProcessor) GetToken(ctx context.Context, tokenString string) (token *Token, err error) {
	rawClaim := Claim{}
	tk, err := jwt.ParseWithClaims(tokenString, &rawClaim, func(token *jwt.Token) (interface{}, error) {
		return p.jwtSecret, nil
	})

	if err != nil || !tk.Valid {
		return nil, errors.New("invalid token")
	}

	claims := &models.JWTClaim{}
	_ = proto.Unmarshal(rawClaim.UserInfo, claims)
	if claims.UserId < 1 {
		return nil, common.UnauthorizedError
	}

	return &Token{Token: tk}, nil
}

func (p *TokenProcessor) ExtractMetadata(token *Token) map[string]string {
	metadata := make(map[string]string)
	// Extract the claims
	claim := GetJWTClaimFromToken(token)
	if claim.UserId < 1 {
		return metadata
	}

	for _, header := range common.ExtraDataHeaders {
		switch header {
		case common.HeaderUserId:
			metadata[header] = cast.ToString(claim.UserId)
		case common.HeaderChildUserIds:
			metadata[header] = utils.ToJSONString(claim.SubUserIds)
		}
	}

	return metadata
}

func (p *TokenProcessor) CheckPermissions(token *Token, requirePermissions map[int64]int64) bool {
	if len(requirePermissions) == 0 {
		return true
	}

	claim := GetJWTClaimFromToken(token)
	if claim.UserId < 1 {
		return false
	}

	for permissionGroupId, permissionSum := range requirePermissions {
		if !permissions.HasPermission(claim.Permissions[permissionGroupId], permissionSum) {
			return false
		}
	}

	return true
}

func (p *TokenProcessor) GenerateToken(ctx context.Context, user *models.User) (string, error) {
	// Create the JWT claim structure and set the required permissions
	exp := config.ViperGetDurationWithDefault("jwt.exp", time.Hour*24*7)

	userInfo := &models.JWTClaim{
		UserId:    user.Id,
		Signature: user.GetSessionId(), // using session_id as signature, todo using other fields
	}

	userInfoBytes, _ := proto.Marshal(userInfo)
	claim := Claim{
		UserInfo: userInfoBytes,
		MapClaims: jwt.MapClaims{
			"exp": jwt.NewNumericDate(time.Now().Add(exp)), // 7 days
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(p.jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetJWTClaimFromToken(t *Token) *models.JWTClaim {
	rawClaim, ok := t.Claims.(*Claim)
	if !ok {
		return &models.JWTClaim{}
	}

	c := &models.JWTClaim{}
	_ = proto.Unmarshal(rawClaim.UserInfo, c)
	return c
}
