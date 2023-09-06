// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"golang.org/x/exp/slog"
	"reflect"
)

type BasicCallMaker struct {
	E       Emitter
	Actions map[reflect.Type]string
}

func (b BasicCallMaker) Send(ctx context.Context, chargeStationId string, request ocpp.Request) error {
	action, ok := b.Actions[reflect.TypeOf(request)]
	if !ok {
		return fmt.Errorf("unknown request type: %T", request)
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg := &Message{
		MessageType:    MessageTypeCall,
		MessageId:      uuid.New().String(),
		Action:         action,
		RequestPayload: requestBytes,
	}

	slog.Info("sending message", "action", msg.Action, "chargeStationId", chargeStationId)
	return b.E.Emit(ctx, chargeStationId, msg)
}

type DataTransferAction struct {
	VendorId  string
	MessageId string
}

type DataTransferCallMaker struct {
	E       Emitter
	Actions map[reflect.Type]DataTransferAction
}

func (d DataTransferCallMaker) Send(ctx context.Context, chargeStationId string, request ocpp.Request) error {
	dta, ok := d.Actions[reflect.TypeOf(request)]
	if !ok {
		return fmt.Errorf("unknown request type: %T", request)
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}
	requestBytesStr := string(requestBytes)

	dataTransferRequest := ocpp16.DataTransferJson{
		VendorId:  dta.VendorId,
		MessageId: &dta.MessageId,
		Data:      &requestBytesStr,
	}

	dataTransferBytes, err := json.Marshal(dataTransferRequest)
	if err != nil {
		return fmt.Errorf("marshaling data transfer request: %w", err)
	}

	msg := &Message{
		MessageType:    MessageTypeCall,
		MessageId:      uuid.New().String(),
		Action:         "DataTransfer",
		RequestPayload: dataTransferBytes,
	}

	return d.E.Emit(ctx, chargeStationId, msg)
}
