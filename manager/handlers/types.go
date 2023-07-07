package handlers

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
)

type CallHandler interface {
	HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error)
}

type CallHandlerFunc func(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error)

func (ch CallHandlerFunc) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	return ch(ctx, chargeStationId, request)
}

type CallRoute struct {
	NewRequest     func() ocpp.Request
	RequestSchema  string
	ResponseSchema string
	Handler        CallHandler
}

type CallResultHandler interface {
	HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error
}

type CallResultHandlerFunc func(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error

func (crh CallResultHandlerFunc) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	return crh(ctx, chargeStationId, request, response, state)
}

type CallResultRoute struct {
	NewRequest     func() ocpp.Request
	NewResponse    func() ocpp.Response
	RequestSchema  string
	ResponseSchema string
	Handler        CallResultHandler
}

type CallMaker interface {
	Send(ctx context.Context, chargeStationId string, request ocpp.Request) error
}
