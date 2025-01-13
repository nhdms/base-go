package handlers

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/nhdms/base-go/cmd/services/user-service/tables"
	"github.com/nhdms/base-go/pkg/dbtool"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/nhdms/base-go/proto/exmsg/services"
)

type UserHandler struct {
	db    *dbtool.ConnectionManager
	redis *redis.Client
}

func (u *UserHandler) GetUserByID(ctx context.Context, request *services.UserRequest, response *services.UserResponse) error {
	sqlTool := dbtool.NewSelect(ctx, u.db.GetConnection(), tables.GetUserTable(), &models.User{})

	qb := squirrel.
		Select(sqlTool.GetQueryColumnList("u")...).
		From(sqlTool.GetTable("u")).
		Where(squirrel.Eq{"id": request.UserId}).
		Limit(1)

	response.User = &models.User{}
	err := sqlTool.Get(ctx, response.User, qb)
	if err != nil {
		logger.DefaultLogger.Errorw("Failed to scan row", "error", err)
		return err
	}
	return nil
}

func NewUserHandler(db *dbtool.ConnectionManager, redis *redis.Client) *UserHandler {
	return &UserHandler{db: db, redis: redis}
}
