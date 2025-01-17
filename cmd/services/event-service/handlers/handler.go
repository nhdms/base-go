package handlers

import (
	"context"
	"gitlab.com/a7923/athena-go/pkg/dbtool"
	"gitlab.com/a7923/athena-go/proto/exmsg/services"
)

type EventHandler struct {
	db *dbtool.ConnectionManager
}

func (e EventHandler) GetEvents(ctx context.Context, request *services.EventRequest, response *services.EventResponse) error {
	//TODO implement me
	panic("implement me")
}

func NewEventHandler(db *dbtool.ConnectionManager) *EventHandler {
	return &EventHandler{db: db}
}
