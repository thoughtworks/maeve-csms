// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/google/uuid"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"golang.org/x/exp/slog"
)

// OcppCallMaker is an implementation of the CallMaker interface for a specific set of OCPP messages.
type OcppCallMaker struct {
	Emitter     transport.Emitter       // used to send the message to the charge station
	OcppVersion transport.OcppVersion   // identifies the OCPP version that the messages are for
	Actions     map[reflect.Type]string // the OCPP Action associated with a specific ocpp.Request object
}

type SetChargingProfileRequestJsonFix struct {
	evseId          int
	chargingProfile interface{}
}

func (b OcppCallMaker) Send(ctx context.Context, chargeStationId string, request ocpp.Request) error {
	action, ok := b.Actions[reflect.TypeOf(request)]
	slog.Info("[TEST] we are in Send() in call_maker.go", "action", action)
	if !ok {
		slog.Error("unknown request type", request)
		return nil
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg := &transport.Message{
		MessageType:    transport.MessageTypeCall,
		MessageId:      uuid.New().String(),
		Action:         action,
		RequestPayload: requestBytes,
	}

	slog.Info("sending message", "action", msg.Action, "chargeStationId", chargeStationId)
	return b.Emitter.Emit(ctx, b.OcppVersion, chargeStationId, msg)
}
