// SPDX-License-Identifier: Apache-2.0

package transport

import (
	"context"
)

// OcppVersion represents the version of OCPP that is being used.
type OcppVersion string

const (
	OcppVersion16  OcppVersion = "ocpp1.6"   // OCPP 1.6
	OcppVersion201 OcppVersion = "ocpp2.0.1" // OCPP 2.0.1
)

// Emitter defines the contract for sending messages to the gateway.
type Emitter interface {
	// Emit sends a message, destined for a specific charge station which is identified by its
	// chargeStationId using a specific ocppVersion, to the gateway.
	Emit(ctx context.Context, ocppVersion OcppVersion, chargeStationId string, message *Message) error
}

// EmitterFunc allows a plain function to be used as an Emitter
type EmitterFunc func(ctx context.Context, ocppVersion OcppVersion, chargeStationId string, message *Message) error

func (e EmitterFunc) Emit(ctx context.Context, ocppVersion OcppVersion, chargeStationId string, message *Message) error {
	return e(ctx, ocppVersion, chargeStationId, message)
}
