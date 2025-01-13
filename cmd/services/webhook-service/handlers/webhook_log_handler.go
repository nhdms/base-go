package handlers

import (
	"context"
	"github.com/Masterminds/squirrel"
	tables2 "github.com/nhdms/base-go/cmd/services/webhook-service/tables"
	"github.com/nhdms/base-go/pkg/dbtool"
	"github.com/nhdms/base-go/pkg/logger"
	"github.com/nhdms/base-go/proto/exmsg/models"
)

type WebhookHandler struct {
	db *dbtool.ConnectionManager
}

func NewWebhookHandler(db *dbtool.ConnectionManager) *WebhookHandler {
	return &WebhookHandler{db: db}
}

func (w *WebhookHandler) InsertLogs(ctx context.Context, events *models.WebhookEvents, result *models.SQLResult) error {
	sqlTool := dbtool.NewInsert(ctx, w.db.GetConnection(), tables2.GetWebhookEventsTable(), &models.WebhookEvent{})
	qb := squirrel.
		Insert(sqlTool.GetTable("")).
		Columns(sqlTool.GetQueryColumnList("")...)

	for _, event := range events.Events {
		qb = qb.Values(sqlTool.GetFilledValues(event)...)
	}

	qb = qb.Suffix("ON CONFLICT (message_uuid) DO UPDATE SET updated_at = now()")

	var err error
	result, err = sqlTool.Insert(ctx, qb)
	if err != nil {
		logger.DefaultLogger.Errorw("Failed to insert webhook logs", "error", err)
		return err
	}
	return nil
}
