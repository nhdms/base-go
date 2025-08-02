package handlers

import (
	"context"
	"github.com/nhdms/base-go/pkg/dbtool"
	"github.com/nhdms/base-go/proto/exmsg/services"
)

type SampleHandler struct {
	db *dbtool.ConnectionManager
}

func (e *SampleHandler) GetSamples(ctx context.Context, request *services.SampleRequest, response *services.SampleResponse) error {
	//TODO implement me
	panic("implement me")
}

func NewSampleHandler(db *dbtool.ConnectionManager) *SampleHandler {
	return &SampleHandler{db: db}
}
