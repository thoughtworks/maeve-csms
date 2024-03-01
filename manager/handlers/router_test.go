// SPDX-License-Identifier: Apache-2.0

package handlers_test

import (
	"context"
	"errors"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"k8s.io/utils/clock"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
)

var heartbeatMsg = transport.Message{
	Action:         "Heartbeat",
	MessageType:    transport.MessageTypeCall,
	RequestPayload: []byte("{}"),
}

var resultMsg = transport.Message{
	Action:          "Result",
	MessageType:     transport.MessageTypeCallResult,
	RequestPayload:  []byte("{}"),
	ResponsePayload: []byte("{}"),
}

var fakeEmitter = func(ctx context.Context, ocppVersion transport.OcppVersion, id string, msg *transport.Message) error {
	return nil
}

func TestRouterHandlesCall(t *testing.T) {
	var chargeStationId string
	var message *transport.Message

	emitter := func(ctx context.Context, ocppVersion transport.OcppVersion, id string, msg *transport.Message) error {
		chargeStationId = id
		message = msg
		return nil
	}

	router := handlers.Router{
		Emitter:     transport.EmitterFunc(emitter),
		SchemaFS:    schemas.OcppSchemas,
		OcppVersion: transport.OcppVersion201,
		CallRoutes: map[string]handlers.CallRoute{
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				RequestSchema:  "ocpp201/HeartbeatRequest.json",
				ResponseSchema: "ocpp201/HeartbeatResponse.json",
				Handler: handlers201.HeartbeatHandler{
					Clock: clock.RealClock{},
				},
			},
		},
	}

	err := router.Route(context.Background(), "id", heartbeatMsg)

	assert.NoError(t, err)
	assert.Equal(t, "id", chargeStationId)
	assert.Equal(t, transport.MessageTypeCallResult, message.MessageType)
	assert.Equal(t, "Heartbeat", message.Action)
	assert.Equal(t, "", message.MessageId)
	assert.NotNil(t, message.ResponsePayload)
}

func TestRouterErrorWhenNoCallRoute(t *testing.T) {
	router := handlers.Router{
		Emitter:     transport.EmitterFunc(fakeEmitter),
		SchemaFS:    schemas.OcppSchemas,
		OcppVersion: transport.OcppVersion201,
		CallRoutes:  map[string]handlers.CallRoute{},
	}

	err := router.Route(context.Background(), "id", heartbeatMsg)

	var mqttError *transport.Error
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mqttError)
	assert.Equal(t, transport.ErrorNotImplemented, mqttError.ErrorCode)
}

func TestRouterErrorWhenCallRequestPayloadIsInvalid(t *testing.T) {
	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: schemas.OcppSchemas,
		CallRoutes: map[string]handlers.CallRoute{
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(noUnmarshalRequest) },
				RequestSchema:  "ocpp201/Heartbeat.json",
				ResponseSchema: "ocpp201/Heartbeat.json",
				Handler: handlers201.HeartbeatHandler{
					Clock: clock.RealClock{},
				},
			},
		},
	}

	var heartbeatMsgWithEmptyPayload = transport.Message{
		Action:         "Heartbeat",
		MessageType:    transport.MessageTypeCall,
		RequestPayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", heartbeatMsgWithEmptyPayload)

	assert.ErrorContains(t, err, "validating Heartbeat request")
}

func TestRouterErrorWhenCantUnmarshallCallRequestPayload(t *testing.T) {
	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: os.DirFS("testdata"),
		CallRoutes: map[string]handlers.CallRoute{
			"MyCall": {
				NewRequest:     func() ocpp.Request { return new(noUnmarshalRequest) },
				RequestSchema:  "schemas/EmptySchema.json",
				ResponseSchema: "schemas/EmptySchema.json",
				Handler: handlers201.HeartbeatHandler{
					Clock: clock.RealClock{},
				},
			},
		},
	}

	var heartbeatMsgWithEmptyPayload = transport.Message{
		Action:         "MyCall",
		MessageType:    transport.MessageTypeCall,
		RequestPayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", heartbeatMsgWithEmptyPayload)

	assert.ErrorContains(t, err, "unmarshalling MyCall request payload")
}

func TestRouterErrorWhenCallHandlerErrors(t *testing.T) {
	handler := func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
		return nil, errors.New("handler error")
	}

	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: schemas.OcppSchemas,
		CallRoutes: map[string]handlers.CallRoute{
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				RequestSchema:  "ocpp201/HeartbeatRequest.json",
				ResponseSchema: "ocpp201/HeartbeatResponse.json",
				Handler:        handlers.CallHandlerFunc(handler),
			},
		},
	}

	err := router.Route(context.Background(), "id", heartbeatMsg)

	assert.ErrorContains(t, err, "handler error")
}

func TestRouterErrorWhenErrorMarshallingCallHandlerResponse(t *testing.T) {
	handler := func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
		return new(noMarshalResponse), nil
	}

	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: schemas.OcppSchemas,
		CallRoutes: map[string]handlers.CallRoute{
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				RequestSchema:  "ocpp201/HeartbeatRequest.json",
				ResponseSchema: "ocpp201/HeartbeatResponse.json",
				Handler:        handlers.CallHandlerFunc(handler),
			},
		},
	}

	err := router.Route(context.Background(), "id", heartbeatMsg)

	assert.ErrorContains(t, err, "marshalling Heartbeat call response")
}

func TestRouterHandlesCallResult(t *testing.T) {
	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: os.DirFS("testdata"),
		CallResultRoutes: map[string]handlers.CallResultRoute{
			"Result": {
				NewRequest:     func() ocpp.Request { return new(fakeRequest) },
				NewResponse:    func() ocpp.Response { return new(fakeResponse) },
				RequestSchema:  "schemas/EmptySchema.json",
				ResponseSchema: "schemas/EmptySchema.json",
				Handler:        handlers.CallResultHandlerFunc(handler),
			},
		},
	}

	err := router.Route(context.Background(), "id", resultMsg)

	assert.NoError(t, err)
}

func TestRouterErrorWhenNoCallResultRoute(t *testing.T) {
	router := handlers.Router{
		Emitter:          transport.EmitterFunc(fakeEmitter),
		SchemaFS:         os.DirFS("testdata"),
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}

	err := router.Route(context.Background(), "id", resultMsg)

	var mqttError *transport.Error
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mqttError)
	assert.Equal(t, transport.ErrorNotImplemented, mqttError.ErrorCode)
}

func TestRouterErrorWhenInvalidCallResultRequestPayload(t *testing.T) {
	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: os.DirFS("testdata"),
		CallResultRoutes: map[string]handlers.CallResultRoute{
			"MyCallResult": {
				NewRequest:     func() ocpp.Request { return nil },
				NewResponse:    func() ocpp.Response { return nil },
				RequestSchema:  "schemas/RequiredFieldSchema.json",
				ResponseSchema: "schemas/EmptySchema.json",
				Handler:        handlers.CallResultHandlerFunc(handler),
			},
		},
	}

	var resultWithEmptyRequestPayload = transport.Message{
		Action:          "MyCallResult",
		MessageType:     transport.MessageTypeCallResult,
		RequestPayload:  []byte("{}"),
		ResponsePayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", resultWithEmptyRequestPayload)

	assert.ErrorContains(t, err, "validating MyCallResult request")
}

func TestRouterErrorWhenCantUnmarshallCallResultRequestPayload(t *testing.T) {
	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: os.DirFS("testdata"),
		CallResultRoutes: map[string]handlers.CallResultRoute{
			"MyCallResult": {
				NewRequest:     func() ocpp.Request { return nil },
				NewResponse:    func() ocpp.Response { return nil },
				RequestSchema:  "schemas/EmptySchema.json",
				ResponseSchema: "schemas/EmptySchema.json",
				Handler:        handlers.CallResultHandlerFunc(handler),
			},
		},
	}

	var resultWithEmptyRequestPayload = transport.Message{
		Action:          "MyCallResult",
		MessageType:     transport.MessageTypeCallResult,
		RequestPayload:  []byte("{}"),
		ResponsePayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", resultWithEmptyRequestPayload)

	assert.ErrorContains(t, err, "unmarshalling MyCallResult request payload")
}

func TestRouterErrorWhenInvalidCallResultResponsePayload(t *testing.T) {
	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: os.DirFS("testdata"),
		CallResultRoutes: map[string]handlers.CallResultRoute{
			"MyCallResult": {
				NewRequest:     func() ocpp.Request { return nil },
				NewResponse:    func() ocpp.Response { return nil },
				RequestSchema:  "schemas/EmptySchema.json",
				ResponseSchema: "schemas/RequiredFieldSchema.json",
				Handler:        handlers.CallResultHandlerFunc(handler),
			},
		},
	}

	var resultWithEmptyRequestPayload = transport.Message{
		Action:          "MyCallResult",
		MessageType:     transport.MessageTypeCallResult,
		RequestPayload:  []byte("{}"),
		ResponsePayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", resultWithEmptyRequestPayload)

	assert.ErrorContains(t, err, "validating MyCallResult response")
}

func TestRouterErrorWhenCantUnmarshallCallResultResponsePayload(t *testing.T) {
	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: os.DirFS("testdata"),
		CallResultRoutes: map[string]handlers.CallResultRoute{
			"MyCallResult": {
				NewRequest:     func() ocpp.Request { return new(fakeRequest) },
				NewResponse:    func() ocpp.Response { return new(noUnmarshalResponse) },
				RequestSchema:  "schemas/EmptySchema.json",
				ResponseSchema: "schemas/EmptySchema.json",
				Handler:        handlers.CallResultHandlerFunc(handler),
			},
		},
	}

	var resultWithEmptyResponsePayload = transport.Message{
		Action:          "MyCallResult",
		MessageType:     transport.MessageTypeCallResult,
		RequestPayload:  []byte("{}"),
		ResponsePayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", resultWithEmptyResponsePayload)

	assert.ErrorContains(t, err, "unmarshalling MyCallResult response payload")
}

func TestRouterErrorWhenCallResultHandlerErrors(t *testing.T) {
	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return errors.New("handler error")
	}

	router := handlers.Router{
		Emitter:  transport.EmitterFunc(fakeEmitter),
		SchemaFS: os.DirFS("testdata"),
		CallResultRoutes: map[string]handlers.CallResultRoute{
			"Result": {
				NewRequest:     func() ocpp.Request { return new(fakeRequest) },
				NewResponse:    func() ocpp.Response { return new(fakeResponse) },
				RequestSchema:  "schemas/EmptySchema.json",
				ResponseSchema: "schemas/EmptySchema.json",
				Handler:        handlers.CallResultHandlerFunc(handler),
			},
		},
	}

	err := router.Route(context.Background(), "id", resultMsg)

	assert.ErrorContains(t, err, "handler error")
}

type fakeRequest struct{}

func (*fakeRequest) IsRequest() {}

type fakeResponse struct{}

func (*fakeResponse) IsResponse() {}

type noUnmarshalRequest struct{}

func (*noUnmarshalRequest) IsRequest() {}

func (*noUnmarshalRequest) UnmarshalJSON(data []byte) error {
	return errors.New("expected to fail")
}

type noUnmarshalResponse struct{}

func (*noUnmarshalResponse) IsResponse() {}

func (*noUnmarshalResponse) UnmarshalJSON(data []byte) error {
	return errors.New("expected to fail")
}

type noMarshalResponse struct{}

func (*noMarshalResponse) IsResponse() {}

func (*noMarshalResponse) MarshalJSON() ([]byte, error) {
	return nil, errors.New("expected to fail")
}
