package mqtt_test

import (
	"context"
	"errors"
	"github.com/twlabs/maeve-csms/manager/handlers"
	handlers201 "github.com/twlabs/maeve-csms/manager/handlers/ocpp201"
	"github.com/twlabs/maeve-csms/manager/schemas"
	"k8s.io/utils/clock"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twlabs/maeve-csms/manager/mqtt"
	"github.com/twlabs/maeve-csms/manager/ocpp"
	"github.com/twlabs/maeve-csms/manager/ocpp/ocpp201"
)

var heartbeatMsg = mqtt.Message{
	Action:         "Heartbeat",
	MessageType:    mqtt.MessageTypeCall,
	RequestPayload: []byte("{}"),
}

var resultMsg = mqtt.Message{
	Action:          "Result",
	MessageType:     mqtt.MessageTypeCallResult,
	RequestPayload:  []byte("{}"),
	ResponsePayload: []byte("{}"),
}

var fakeEmitter = func(ctx context.Context, id string, msg *mqtt.Message) error {
	return nil
}

func TestRouterHandlesCall(t *testing.T) {
	router := mqtt.Router{
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

	var chargeStationId string
	var message *mqtt.Message
	emitter := func(ctx context.Context, id string, msg *mqtt.Message) error {
		chargeStationId = id
		message = msg
		return nil
	}

	err := router.Route(context.Background(), "id", heartbeatMsg, mqtt.EmitterFunc(emitter), schemas.OcppSchemas)

	assert.NoError(t, err)
	assert.Equal(t, "id", chargeStationId)
	assert.Equal(t, mqtt.MessageTypeCallResult, message.MessageType)
	assert.Equal(t, "Heartbeat", message.Action)
	assert.Equal(t, "", message.MessageId)
	assert.NotNil(t, message.ResponsePayload)
}

func TestRouterErrorWhenNoCallRoute(t *testing.T) {
	router := mqtt.Router{
		CallRoutes: map[string]handlers.CallRoute{},
	}

	err := router.Route(context.Background(), "id", heartbeatMsg, mqtt.EmitterFunc(fakeEmitter), schemas.OcppSchemas)

	var mqttError *mqtt.Error
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mqttError)
	assert.Equal(t, mqtt.ErrorNotImplemented, mqttError.ErrorCode)
}

func TestRouterErrorWhenCallRequestPayloadIsInvalid(t *testing.T) {
	router := mqtt.Router{
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

	var heartbeatMsgWithEmptyPayload = mqtt.Message{
		Action:         "Heartbeat",
		MessageType:    mqtt.MessageTypeCall,
		RequestPayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", heartbeatMsgWithEmptyPayload, mqtt.EmitterFunc(fakeEmitter), schemas.OcppSchemas)

	assert.ErrorContains(t, err, "validating Heartbeat request")
}

func TestRouterErrorWhenCantUnmarshallCallRequestPayload(t *testing.T) {
	router := mqtt.Router{
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

	var heartbeatMsgWithEmptyPayload = mqtt.Message{
		Action:         "MyCall",
		MessageType:    mqtt.MessageTypeCall,
		RequestPayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", heartbeatMsgWithEmptyPayload, mqtt.EmitterFunc(fakeEmitter), os.DirFS("testdata"))

	assert.ErrorContains(t, err, "unmarshalling MyCall request payload")
}

func TestRouterErrorWhenCallHandlerErrors(t *testing.T) {
	handler := func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
		return nil, errors.New("handler error")
	}

	router := mqtt.Router{
		CallRoutes: map[string]handlers.CallRoute{
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				RequestSchema:  "ocpp201/HeartbeatRequest.json",
				ResponseSchema: "ocpp201/HeartbeatResponse.json",
				Handler:        handlers.CallHandlerFunc(handler),
			},
		},
	}

	err := router.Route(context.Background(), "id", heartbeatMsg, mqtt.EmitterFunc(fakeEmitter), schemas.OcppSchemas)

	assert.ErrorContains(t, err, "handler error")
}

func TestRouterErrorWhenErrorMarshallingCallHandlerResponse(t *testing.T) {
	handler := func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
		return new(noMarshalResponse), nil
	}

	router := mqtt.Router{
		CallRoutes: map[string]handlers.CallRoute{
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				RequestSchema:  "ocpp201/HeartbeatRequest.json",
				ResponseSchema: "ocpp201/HeartbeatResponse.json",
				Handler:        handlers.CallHandlerFunc(handler),
			},
		},
	}

	err := router.Route(context.Background(), "id", heartbeatMsg, mqtt.EmitterFunc(fakeEmitter), schemas.OcppSchemas)

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

	router := mqtt.Router{
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

	err := router.Route(context.Background(), "id", resultMsg, mqtt.EmitterFunc(fakeEmitter), os.DirFS("testdata"))

	assert.NoError(t, err)
}

func TestRouterErrorWhenNoCallResultRoute(t *testing.T) {
	router := mqtt.Router{
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}

	err := router.Route(context.Background(), "id", resultMsg, mqtt.EmitterFunc(fakeEmitter), os.DirFS("testdata"))

	var mqttError *mqtt.Error
	assert.Error(t, err)
	assert.ErrorAs(t, err, &mqttError)
	assert.Equal(t, mqtt.ErrorNotImplemented, mqttError.ErrorCode)
}

func TestRouterErrorWhenInvalidCallResultRequestPayload(t *testing.T) {
	handler := func(ctx context.Context,
		chargeStationId string,
		request ocpp.Request,
		response ocpp.Response,
		state any) error {
		return nil
	}

	router := mqtt.Router{
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

	var resultWithEmptyRequestPayload = mqtt.Message{
		Action:          "MyCallResult",
		MessageType:     mqtt.MessageTypeCallResult,
		RequestPayload:  []byte("{}"),
		ResponsePayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", resultWithEmptyRequestPayload, mqtt.EmitterFunc(fakeEmitter), os.DirFS("testdata"))

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

	router := mqtt.Router{
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

	var resultWithEmptyRequestPayload = mqtt.Message{
		Action:          "MyCallResult",
		MessageType:     mqtt.MessageTypeCallResult,
		RequestPayload:  []byte("{}"),
		ResponsePayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", resultWithEmptyRequestPayload, mqtt.EmitterFunc(fakeEmitter), os.DirFS("testdata"))

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

	router := mqtt.Router{
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

	var resultWithEmptyRequestPayload = mqtt.Message{
		Action:          "MyCallResult",
		MessageType:     mqtt.MessageTypeCallResult,
		RequestPayload:  []byte("{}"),
		ResponsePayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", resultWithEmptyRequestPayload, mqtt.EmitterFunc(fakeEmitter), os.DirFS("testdata"))

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

	router := mqtt.Router{
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

	var resultWithEmptyResponsePayload = mqtt.Message{
		Action:          "MyCallResult",
		MessageType:     mqtt.MessageTypeCallResult,
		RequestPayload:  []byte("{}"),
		ResponsePayload: []byte("{}"),
	}

	err := router.Route(context.Background(), "id", resultWithEmptyResponsePayload, mqtt.EmitterFunc(fakeEmitter), os.DirFS("testdata"))

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

	router := mqtt.Router{
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

	err := router.Route(context.Background(), "id", resultMsg, mqtt.EmitterFunc(fakeEmitter), os.DirFS("testdata"))

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
