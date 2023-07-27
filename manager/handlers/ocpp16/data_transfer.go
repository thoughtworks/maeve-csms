// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io/fs"

	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"golang.org/x/exp/slog"
)

type DataTransferHandler struct {
	CallRoutes map[string]map[string]handlers.CallRoute
	SchemaFS   fs.FS
}

func (d DataTransferHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	span := trace.SpanFromContext(ctx)

	req := request.(*types.DataTransferJson)

	messageId := ""
	if req.MessageId != nil {
		messageId = *req.MessageId
	}
	slog.Info("data transfer",
		slog.String("vendorId", req.VendorId), slog.String("messageId", messageId))

	span.SetAttributes(attribute.String("datatransfer.vendor_id", req.VendorId))
	if messageId != "" {
		span.SetAttributes(attribute.String("datatransfer.message_id", messageId))
	}

	vendorMap, ok := d.CallRoutes[req.VendorId]
	if !ok {
		return &types.DataTransferResponseJson{
			Status: types.DataTransferResponseJsonStatusUnknownVendorId,
		}, nil
	}
	route, ok := vendorMap[messageId]
	if !ok {
		return &types.DataTransferResponseJson{
			Status: types.DataTransferResponseJsonStatusUnknownMessageId,
		}, nil
	}

	var dataTransferRequest ocpp.Request
	if req.Data != nil {
		data := []byte(*req.Data)
		err := schemas.Validate(data, d.SchemaFS, route.RequestSchema)
		if err != nil {
			return nil, fmt.Errorf("validating %s:%s data transfer data: %w", req.VendorId, messageId, err)
		}
		dataTransferRequest = route.NewRequest()
		err = json.Unmarshal(data, &dataTransferRequest)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling %s:%s data transfer data: %w", req.VendorId, messageId, err)
		}
	}

	dataTransferResponse, err := route.Handler.HandleCall(ctx, chargeStationId, dataTransferRequest)
	if err != nil {
		return nil, err
	}
	var dataTransferResponseData *string
	if dataTransferResponse != nil {
		b, err := json.Marshal(dataTransferResponse)
		if err != nil {
			return nil, fmt.Errorf("marshalling %s:%s data transfer data: %w", req.VendorId, messageId, err)
		}
		err = schemas.Validate(b, d.SchemaFS, route.ResponseSchema)
		if err != nil {
			slog.Warn("data transfer response is not valid", slog.String("vendorId", req.VendorId), slog.String("messageId", messageId), err)
		}
		dataTransferResponseString := string(b)
		dataTransferResponseData = &dataTransferResponseString
	}

	return &types.DataTransferResponseJson{
		Status: types.DataTransferResponseJsonStatusAccepted,
		Data:   dataTransferResponseData,
	}, nil
}
