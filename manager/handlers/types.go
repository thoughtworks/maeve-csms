// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
)

// CallHandler is the interface implemented by handlers that are designed to process an OCPP Call.
type CallHandler interface {
	// HandleCall receives the charge station identifier and the OCPP request message
	// and returns either the OCPP response message or an error.
	HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error)
}

// CallHandlerFunc allows a plain function to be used as a CallHandler.
type CallHandlerFunc func(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error)

func (ch CallHandlerFunc) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	return ch(ctx, chargeStationId, request)
}

// CallRoute is the configuration that is used by the Router for processing an OCPP Call.
// In the Router this is indexed by the OCPP Action.
type CallRoute struct {
	NewRequest     func() ocpp.Request // Function used for creating an empty request
	RequestSchema  string              // JSON schema file that corresponds to the request data structure
	ResponseSchema string              // JSON schema file that corresponds to the response data structure
	Handler        CallHandler         // Function to process a call
}

// CallResultHandler is the interface implemented by the handlers that are designed to process an OCPP CallResult.
type CallResultHandler interface {
	// HandleCallResult receives the charge station id, OCPP Request message and OCPP Response message
	// along with any cached state. It may return an error.
	HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error
}

// CallResultHandlerFunc allows a plain function to be used as a CallResultHandler
type CallResultHandlerFunc func(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error

func (crh CallResultHandlerFunc) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	return crh(ctx, chargeStationId, request, response, state)
}

// CallResultRoute is the configuration that is used by the Router for process an OCPP CallResult.
// In the Router this is indexed by the OCPP Action (of the corresponding Call).
type CallResultRoute struct {
	NewRequest     func() ocpp.Request  // Function used for creating an empty request
	NewResponse    func() ocpp.Response // Function used for creating an empty response
	RequestSchema  string               // JSON schema file that corresponds to the request data structure
	ResponseSchema string               // JSON schema file that corresponds to the response data structure
	Handler        CallResultHandler    // Function to process a call result
}

// CallMaker is the interface used by handlers (and other parts of the system) that want to initiate
// an OCPP call from the CSMS.
type CallMaker interface {
	// Send receives the charge station id and the request to send. It may return an error.
	Send(ctx context.Context, chargeStationId string, request ocpp.Request) error
}
