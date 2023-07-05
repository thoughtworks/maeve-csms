package ocpp16

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/twlabs/maeve-csms/manager/handlers"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp16"
	"github.com/twlabs/maeve-csms/manager/schemas"
	"io/fs"
	"log"
)

type DataTransferHandler struct {
	CallRoutes map[string]map[string]handlers.CallRoute
	SchemaFS   fs.FS
}

func (d DataTransferHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.DataTransferJson)

	messageId := ""
	if req.MessageId != nil {
		messageId = *req.MessageId
	}
	log.Printf("Data transfer %s:%s", req.VendorId, messageId)

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
			log.Printf("warning: data transfer response to %s:%s is not valid: %v", req.VendorId, messageId, err)
		}
		dataTransferResponseString := string(b)
		dataTransferResponseData = &dataTransferResponseString
	}

	return &types.DataTransferResponseJson{
		Status: types.DataTransferResponseJsonStatusAccepted,
		Data:   dataTransferResponseData,
	}, nil
}
