// SPDX-License-Identifier: Apache-2.0

package handlers_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"github.com/thoughtworks/maeve-csms/manager/testutil"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"k8s.io/utils/clock"
	"os"
	"strings"
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

func TestRouterHandlesCall(t *testing.T) {
	emitter := new(FakeEmitter)

	router := handlers.Router{
		Emitter:     emitter,
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

	router.Handle(context.Background(), "id", &heartbeatMsg)

	assert.Equal(t, "id", emitter.chargeStationId)
	assert.Equal(t, transport.MessageTypeCallResult, emitter.msg.MessageType)
	assert.Equal(t, "Heartbeat", emitter.msg.Action)
	assert.Equal(t, "", emitter.msg.MessageId)
	assert.NotNil(t, emitter.msg.ResponsePayload)
}

func TestRouterErrorWhenNoCallRoute(t *testing.T) {
	emitter := new(FakeEmitter)

	router := handlers.Router{
		Emitter:     emitter,
		SchemaFS:    schemas.OcppSchemas,
		OcppVersion: transport.OcppVersion201,
		CallRoutes:  map[string]handlers.CallRoute{},
	}

	router.Handle(context.Background(), "id", &heartbeatMsg)

	assert.Equal(t, "id", emitter.chargeStationId)
	assert.Equal(t, transport.MessageTypeCallError, emitter.msg.MessageType)
	assert.Equal(t, "Heartbeat", emitter.msg.Action)
	assert.Equal(t, "", emitter.msg.MessageId)
	assert.Equal(t, transport.ErrorNotImplemented, emitter.msg.ErrorCode)
	assert.Nil(t, emitter.msg.ResponsePayload)
}

func TestRouterErrorWhenCallRequestPayloadIsInvalid(t *testing.T) {
	emitter := new(FakeEmitter)

	router := handlers.Router{
		Emitter:  emitter,
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

	router.Handle(context.Background(), "id", &heartbeatMsgWithEmptyPayload)

	assert.Equal(t, "id", emitter.chargeStationId)
	assert.Equal(t, transport.MessageTypeCallError, emitter.msg.MessageType)
	assert.Equal(t, "Heartbeat", emitter.msg.Action)
	assert.Equal(t, "", emitter.msg.MessageId)
	assert.Equal(t, transport.ErrorFormatViolation, emitter.msg.ErrorCode)
	assert.Nil(t, emitter.msg.ResponsePayload)
}

func TestRouterErrorWhenCantUnmarshallCallRequestPayload(t *testing.T) {
	emitter := new(FakeEmitter)

	router := handlers.Router{
		Emitter:  emitter,
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

	router.Handle(context.Background(), "id", &heartbeatMsgWithEmptyPayload)

	assert.Equal(t, "id", emitter.chargeStationId)
	assert.Equal(t, transport.MessageTypeCallError, emitter.msg.MessageType)
	assert.Equal(t, "MyCall", emitter.msg.Action)
	assert.Equal(t, "", emitter.msg.MessageId)
	assert.Equal(t, transport.ErrorInternalError, emitter.msg.ErrorCode)
	assert.Nil(t, emitter.msg.ResponsePayload)
}

func TestRouterErrorWhenCallHandlerErrors(t *testing.T) {
	emitter := new(FakeEmitter)

	handler := func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
		return nil, errors.New("handler error")
	}

	router := handlers.Router{
		Emitter:  emitter,
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

	router.Handle(context.Background(), "id", &heartbeatMsg)

	assert.Equal(t, "id", emitter.chargeStationId)
	assert.Equal(t, transport.MessageTypeCallError, emitter.msg.MessageType)
	assert.Equal(t, "Heartbeat", emitter.msg.Action)
	assert.Equal(t, "", emitter.msg.MessageId)
	assert.Equal(t, transport.ErrorInternalError, emitter.msg.ErrorCode)
	assert.Nil(t, emitter.msg.ResponsePayload)
}

func TestRouterErrorWhenErrorMarshallingCallHandlerResponse(t *testing.T) {
	emitter := new(FakeEmitter)

	handler := func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
		return new(noMarshalResponse), nil
	}

	router := handlers.Router{
		Emitter:  emitter,
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

	router.Handle(context.Background(), "id", &heartbeatMsg)

	assert.Equal(t, "id", emitter.chargeStationId)
	assert.Equal(t, transport.MessageTypeCallError, emitter.msg.MessageType)
	assert.Equal(t, "Heartbeat", emitter.msg.Action)
	assert.Equal(t, "", emitter.msg.MessageId)
	assert.Equal(t, transport.ErrorInternalError, emitter.msg.ErrorCode)
	assert.Nil(t, emitter.msg.ResponsePayload)
}

func TestRouterHandlesCallResult(t *testing.T) {
	tracer, exporter := testutil.GetTracer()

	emitter := new(FakeEmitter)

	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  emitter,
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

	func() {
		ctx, span := tracer.Start(context.Background(), "test")
		defer span.End()
		router.Handle(ctx, "id", &resultMsg)
	}()

	// for a call response the emitter should never be called
	assert.False(t, emitter.called)

	// check that no error was produced using telemetry
	assert.Greater(t, len(exporter.GetSpans()), 0)
	assert.Equal(t, codes.Ok, exporter.GetSpans()[0].Status.Code)
}

func TestRouterErrorWhenNoCallResultRoute(t *testing.T) {
	tracer, exporter := testutil.GetTracer()

	emitter := new(FakeEmitter)

	router := handlers.Router{
		Emitter:          emitter,
		SchemaFS:         os.DirFS("testdata"),
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}

	func() {
		ctx, span := tracer.Start(context.Background(), "test")
		defer span.End()
		router.Handle(ctx, "id", &resultMsg)
	}()

	// for a call response the emitter should never be called
	assert.False(t, emitter.called)

	// check that an error was produced using telemetry
	require.Greater(t, len(exporter.GetSpans()), 0)
	assert.Equal(t, codes.Error, exporter.GetSpans()[0].Status.Code)
	require.Greater(t, len(exporter.GetSpans()[0].Events), 0)
	testutil.AssertAttributes(t, exporter.GetSpans()[0].Events[0].Attributes, map[string]any{
		"exception.type":    "*fmt.wrapError",
		"exception.message": "routing request: NotImplemented: Result result not implemented",
	})
}

func TestRouterErrorWhenInvalidCallResultRequestPayload(t *testing.T) {
	tracer, exporter := testutil.GetTracer()

	emitter := new(FakeEmitter)

	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  emitter,
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

	func() {
		ctx, span := tracer.Start(context.Background(), "test")
		defer span.End()

		router.Handle(ctx, "id", &resultWithEmptyRequestPayload)
	}()

	// for a call response the emitter should never be called
	assert.False(t, emitter.called)

	// check that an error was produced using telemetry
	require.Greater(t, len(exporter.GetSpans()), 0)
	assert.Equal(t, codes.Error, exporter.GetSpans()[0].Status.Code)
	require.Greater(t, len(exporter.GetSpans()[0].Events), 0)
	testutil.AssertAttributes(t, exporter.GetSpans()[0].Events[0].Attributes, map[string]any{
		"exception.type": "*fmt.wrapError",
		"exception.message": func(val attribute.Value) bool {
			return strings.Contains(val.AsString(), "validating MyCallResult request")
		},
	})
}

func TestRouterErrorWhenCantUnmarshallCallResultRequestPayload(t *testing.T) {
	tracer, exporter := testutil.GetTracer()

	emitter := new(FakeEmitter)

	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  emitter,
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

	func() {
		ctx, span := tracer.Start(context.Background(), "test")
		defer span.End()
		router.Handle(ctx, "id", &resultWithEmptyRequestPayload)
	}()

	// for a call response the emitter should never be called
	assert.False(t, emitter.called)

	// check that an error was produced using telemetry
	require.Greater(t, len(exporter.GetSpans()), 0)
	assert.Equal(t, codes.Error, exporter.GetSpans()[0].Status.Code)
	require.Greater(t, len(exporter.GetSpans()[0].Events), 0)
	testutil.AssertAttributes(t, exporter.GetSpans()[0].Events[0].Attributes, map[string]any{
		"exception.type": "*fmt.wrapError",
		"exception.message": func(val attribute.Value) bool {
			return strings.Contains(val.AsString(), "unmarshalling MyCallResult request payload")
		},
	})
}

func TestRouterErrorWhenInvalidCallResultResponsePayload(t *testing.T) {
	tracer, exporter := testutil.GetTracer()

	emitter := new(FakeEmitter)

	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  emitter,
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

	func() {
		ctx, span := tracer.Start(context.Background(), "test")
		defer span.End()
		router.Handle(ctx, "id", &resultWithEmptyRequestPayload)
	}()

	// for a call response the emitter should never be called
	assert.False(t, emitter.called)

	// check that an error was produced using telemetry
	require.Greater(t, len(exporter.GetSpans()), 0)
	assert.Equal(t, codes.Error, exporter.GetSpans()[0].Status.Code)
	require.Greater(t, len(exporter.GetSpans()[0].Events), 0)
	testutil.AssertAttributes(t, exporter.GetSpans()[0].Events[0].Attributes, map[string]any{
		"exception.type": "*fmt.wrapError",
		"exception.message": func(val attribute.Value) bool {
			return strings.Contains(val.AsString(), "validating MyCallResult response")
		},
	})
}

func TestRouterErrorWhenCantUnmarshallCallResultResponsePayload(t *testing.T) {
	tracer, exporter := testutil.GetTracer()

	emitter := new(FakeEmitter)

	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := handlers.Router{
		Emitter:  emitter,
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

	func() {
		ctx, span := tracer.Start(context.Background(), "test")
		defer span.End()
		router.Handle(ctx, "id", &resultWithEmptyResponsePayload)
	}()

	// for a call response the emitter should never be called
	assert.False(t, emitter.called)

	// check that an error was produced using telemetry
	require.Greater(t, len(exporter.GetSpans()), 0)
	assert.Equal(t, codes.Error, exporter.GetSpans()[0].Status.Code)
	require.Greater(t, len(exporter.GetSpans()[0].Events), 0)
	testutil.AssertAttributes(t, exporter.GetSpans()[0].Events[0].Attributes, map[string]any{
		"exception.type": "*errors.errorString",
		"exception.message": func(val attribute.Value) bool {
			return strings.Contains(val.AsString(), "unmarshalling MyCallResult response payload")
		},
	})
}

func TestRouterErrorWhenCallResultHandlerErrors(t *testing.T) {
	tracer, exporter := testutil.GetTracer()

	emitter := new(FakeEmitter)

	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return errors.New("handler error")
	}

	router := handlers.Router{
		Emitter:  emitter,
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

	func() {
		ctx, span := tracer.Start(context.Background(), "test")
		defer span.End()
		router.Handle(ctx, "id", &resultMsg)
	}()

	// for a call response the emitter should never be called
	assert.False(t, emitter.called)

	// check that an error was produced using telemetry
	require.Greater(t, len(exporter.GetSpans()), 0)
	assert.Equal(t, codes.Error, exporter.GetSpans()[0].Status.Code)
	require.Greater(t, len(exporter.GetSpans()[0].Events), 0)
	testutil.AssertAttributes(t, exporter.GetSpans()[0].Events[0].Attributes, map[string]any{
		"exception.type": "*errors.errorString",
		"exception.message": func(val attribute.Value) bool {
			return strings.Contains(val.AsString(), "handler error")
		},
	})
}

type fakeRequest struct{}

func (*fakeRequest) IsRequest() {}

type fakeResponse struct{}

func (*fakeResponse) IsResponse() {}

type noUnmarshalRequest struct{}

func (*noUnmarshalRequest) IsRequest() {}

func (*noUnmarshalRequest) UnmarshalJSON(_ []byte) error {
	return errors.New("expected to fail")
}

type noUnmarshalResponse struct{}

func (*noUnmarshalResponse) IsResponse() {}

func (*noUnmarshalResponse) UnmarshalJSON(_ []byte) error {
	return errors.New("expected to fail")
}

type noMarshalResponse struct{}

func (*noMarshalResponse) IsResponse() {}

func (*noMarshalResponse) MarshalJSON() ([]byte, error) {
	return nil, errors.New("expected to fail")
}
