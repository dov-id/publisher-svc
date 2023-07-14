package models

import (
	"github.com/dov-id/publisher-svc/internal/data"
	"github.com/dov-id/publisher-svc/resources"
)

func newRequest(request data.Request) resources.Request {
	return resources.Request{
		Key: resources.Key{
			ID:   request.Id,
			Type: resources.REQUEST,
		},
		Attributes: resources.RequestAttributes{
			Id:     request.Id,
			Status: string(request.Status),
			Error:  request.Error,
		},
	}
}

func NewRequestResponse(request data.Request) resources.RequestResponse {
	return resources.RequestResponse{
		Data: newRequest(request),
	}
}
