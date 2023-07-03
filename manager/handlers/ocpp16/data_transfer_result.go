package ocpp16

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/twlabs/ocpp2-broker-core/manager/handlers"
	"github.com/twlabs/ocpp2-broker-core/manager/ocpp"
	types "github.com/twlabs/ocpp2-broker-core/manager/ocpp/ocpp16"
	"github.com/twlabs/ocpp2-broker-core/manager/schemas"
	"io/fs"
	"log"
)

type DataTransferResultHandler struct {
	SchemaFS         fs.FS
	CallResultRoutes map[string]map[string]handlers.CallResultRoute
}

func (d DataTransferResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.DataTransferJson)
	resp := response.(*types.DataTransferResponseJson)

	messageId := ""
	if req.MessageId != nil {
		messageId = *req.MessageId
	}
	log.Printf("Data transfer result %s:%s", req.VendorId, messageId)

	vendorMap, ok := d.CallResultRoutes[req.VendorId]
	if !ok {
		return fmt.Errorf("unknown data transfer result vendor: %s", req.VendorId)
	}
	route, ok := vendorMap[messageId]
	if !ok {
		return fmt.Errorf("unknown data transfer result message id: %s", messageId)
	}

	var dataTransferRequest ocpp.Request
	if req.Data != nil {
		data := []byte(*req.Data)
		err := schemas.Validate(data, d.SchemaFS, route.RequestSchema)
		if err != nil {
			return fmt.Errorf("validating %s:%s data transfer result request data: %w", req.VendorId, messageId, err)
		}
		dataTransferRequest = route.NewRequest()
		err = json.Unmarshal(data, &dataTransferRequest)
		if err != nil {
			return fmt.Errorf("unmarshalling %s:%s data transfer request data: %w", req.VendorId, messageId, err)
		}
	}

	var dataTransferResponse ocpp.Response
	if resp.Data != nil {
		data := []byte(*resp.Data)
		err := schemas.Validate(data, d.SchemaFS, route.ResponseSchema)
		if err != nil {
			return fmt.Errorf("validating %s:%s data transfer result response data: %w", req.VendorId, messageId, err)
		}
		dataTransferResponse = route.NewResponse()
		err = json.Unmarshal(data, &dataTransferResponse)
		if err != nil {
			return fmt.Errorf("unmarshalling %s:%s data transfer response data: %w", req.VendorId, messageId, err)
		}
	}

	return route.Handler.HandleCallResult(ctx, chargeStationId, dataTransferRequest, dataTransferResponse, state)
}
