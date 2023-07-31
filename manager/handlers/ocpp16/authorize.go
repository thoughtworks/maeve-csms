// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
)

type AuthorizeHandler struct {
	TokenStore store.TokenStore
}

func (a AuthorizeHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.AuthorizeJson)
	slog.Info("checking", slog.String("chargeStationId", chargeStationId),
		slog.String("idTag", req.IdTag))

	status := types.AuthorizeResponseJsonIdTagInfoStatusInvalid
	tok, err := a.TokenStore.LookupToken(ctx, req.IdTag)
	if err != nil {
		return nil, err
	}
	if tok != nil {
		status = types.AuthorizeResponseJsonIdTagInfoStatusAccepted
	}

	span.SetAttributes(
		attribute.String("request.status", string(status)),
		attribute.String("authorize.token", req.IdTag))

	return &types.AuthorizeResponseJson{
		IdTagInfo: types.AuthorizeResponseJsonIdTagInfo{
			Status: status,
		},
	}, nil
}
