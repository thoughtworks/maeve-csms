// SPDX-License-Identifier: Apache-2.0

package mqtt_test

import (
	"context"
	"encoding/json"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/mqtt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/utils/clock"
	clockTest "k8s.io/utils/clock/testing"
	"net/url"
	"testing"
	"time"
)

func TestMqttConnection(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start the broker
	broker, clientUrl := mqtt.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()

	err = broker.Serve()
	require.NoError(t, err)

	// connect the mqtt handler
	handler := mqtt.NewHandler(
		mqtt.WithMqttBrokerUrl(clientUrl),
		mqtt.WithMqttPrefix("cs"),
		mqtt.WithClock(clockTest.NewFakePassiveClock(now)))
	errCh := make(chan error)
	handler.Connect(errCh)

	// subscribe to the output channel
	rcvdCh := make(chan struct{})

	mqttClient, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{clientUrl},
		KeepAlive:         10,
		ConnectRetryDelay: 10,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err := manager.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					"cs/out/ocpp2.0.1/cs001": {},
				},
			})
			require.NoError(t, err)
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "test",
			Router: paho.NewSingleHandlerRouter(func(publish *paho.Publish) {
				assert.Equal(t, "cs/out/ocpp2.0.1/cs001", publish.Topic)
				var msg mqtt.Message
				err := json.Unmarshal(publish.Payload, &msg)
				assert.NoError(t, err, "payload is not the expected message type")
				assert.Equal(t, mqtt.MessageTypeCallResult, msg.MessageType)
				assert.Equal(t, "Heartbeat", msg.Action)
				assert.Equal(t, "1234", msg.MessageId)
				assert.JSONEq(t, `{"currentTime":"2023-06-15T15:05:00+01:00"}`, string(msg.ResponsePayload))
				rcvdCh <- struct{}{}
			}),
		},
	})
	require.NoError(t, err)

	err = mqttClient.AwaitConnection(ctx)
	require.NoError(t, err)

	// publish a message to the input channel
	msg := mqtt.Message{
		MessageType:    mqtt.MessageTypeCall,
		Action:         "Heartbeat",
		MessageId:      "1234",
		RequestPayload: []byte("{}"),
	}
	msgBytes, err := json.Marshal(msg)
	require.NoError(t, err)

	err = broker.Publish("cs/in/ocpp2.0.1/cs001", msgBytes, false, 0)
	require.NoError(t, err)

	// wait for success
	select {
	case <-rcvdCh:
		// success
	case err = <-errCh:
		assert.Fail(t, "unexpected error from handler: %v", err)
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test")
	}
}

func TestGatewayMessageHandler(t *testing.T) {
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)

	ctx := context.Background()
	router := &mqtt.Router{
		CallRoutes: map[string]handlers.CallRoute{
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				RequestSchema:  "ocpp201/HeartbeatRequest.json",
				ResponseSchema: "ocpp201/HeartbeatResponse.json",
				Handler: handlers201.HeartbeatHandler{
					Clock: clockTest.NewFakePassiveClock(now),
				},
			},
		},
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}
	var chargeStationId string
	var message *mqtt.Message
	emitter := func(_ context.Context, csId string, msg *mqtt.Message) error {
		chargeStationId = csId
		message = msg
		return nil
	}

	handler := mqtt.NewGatewayMessageHandler(ctx, "test-client", tracer, router, mqtt.EmitterFunc(emitter), schemas.OcppSchemas)

	mqttMsg := mqtt.Message{
		MessageType:    mqtt.MessageTypeCall,
		Action:         "Heartbeat",
		MessageId:      "1234",
		RequestPayload: []byte("{}"),
	}

	mqttMsgBytes, err := json.Marshal(mqttMsg)
	require.NoError(t, err)

	msg := &paho.Publish{
		Topic:   "cs/ocpp2.0.1/cs001",
		Payload: mqttMsgBytes,
	}

	handler(msg)

	assert.Equal(t, "cs001", chargeStationId)
	assert.Equal(t, mqtt.MessageTypeCallResult, message.MessageType)
	assert.Equal(t, "Heartbeat", message.Action)
	assert.Equal(t, "1234", message.MessageId)
	assert.JSONEq(t, `{"currentTime":"2023-06-15T15:05:00+01:00"}`, string(message.ResponsePayload))
}

func TestGatewayMessageHandlerEmitsErrorWhenFailingToUnmarshallIncomingMessage(t *testing.T) {
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	ctx := context.Background()
	router := &mqtt.Router{
		CallRoutes:       map[string]handlers.CallRoute{},
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}
	var chargeStationId string
	var message *mqtt.Message
	emitter := func(_ context.Context, csId string, msg *mqtt.Message) error {
		chargeStationId = csId
		message = msg
		return nil
	}

	handler := mqtt.NewGatewayMessageHandler(ctx, "test-client", tracer, router, mqtt.EmitterFunc(emitter), schemas.OcppSchemas)

	msg := &paho.Publish{
		Topic:   "cs/ocpp2.0.1/cs001",
		Payload: []byte(""),
	}

	handler(msg)

	assert.Equal(t, "cs001", chargeStationId)
	assert.Equal(t, mqtt.MessageTypeCallError, message.MessageType)
	assert.Equal(t, "", message.Action)
	assert.Equal(t, "-1", message.MessageId)
	assert.Equal(t, mqtt.ErrorInternalError, message.ErrorCode)
	assert.NotEmpty(t, message.ErrorDescription)
}

func TestGatewayMessageHandlerEmitsErrorWithCodeFromRouter(t *testing.T) {
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	ctx := context.Background()
	router := &mqtt.Router{
		CallRoutes:       map[string]handlers.CallRoute{},
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}
	var chargeStationId string
	var message *mqtt.Message
	emitter := func(_ context.Context, csId string, msg *mqtt.Message) error {
		chargeStationId = csId
		message = msg
		return nil
	}

	handler := mqtt.NewGatewayMessageHandler(ctx, "test-client", tracer, router, mqtt.EmitterFunc(emitter), schemas.OcppSchemas)

	mqttMsg := mqtt.Message{
		MessageType:    mqtt.MessageTypeCall,
		Action:         "Heartbeat",
		MessageId:      "1234",
		RequestPayload: []byte(""),
	}

	mqttMsgBytes, err := json.Marshal(mqttMsg)
	require.NoError(t, err)

	msg := &paho.Publish{
		Topic:   "cs/ocpp2.0.1/cs001",
		Payload: mqttMsgBytes,
	}

	handler(msg)

	assert.Equal(t, "cs001", chargeStationId)
	assert.Equal(t, mqtt.MessageTypeCallError, message.MessageType)
	assert.Equal(t, "Heartbeat", message.Action)
	assert.Equal(t, "1234", message.MessageId)
	assert.Equal(t, mqtt.ErrorNotImplemented, message.ErrorCode)
	assert.Equal(t, "Heartbeat not implemented", message.ErrorDescription)
}

func TestGatewayMessageHandlerEmitsErrorFromRouter(t *testing.T) {
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	ctx := context.Background()
	router := &mqtt.Router{
		CallRoutes: map[string]handlers.CallRoute{
			"Heartbeat": {
				NewRequest: func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				Handler: handlers201.HeartbeatHandler{
					Clock: clock.RealClock{},
				},
			},
		},
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}
	var chargeStationId string
	var message *mqtt.Message
	emitter := func(_ context.Context, csId string, msg *mqtt.Message) error {
		chargeStationId = csId
		message = msg
		return nil
	}

	handler := mqtt.NewGatewayMessageHandler(ctx, "test-client", tracer, router, mqtt.EmitterFunc(emitter), schemas.OcppSchemas)

	mqttMsg := mqtt.Message{
		MessageType:    mqtt.MessageTypeCall,
		Action:         "Heartbeat",
		MessageId:      "1234",
		RequestPayload: []byte(""),
	}

	mqttMsgBytes, err := json.Marshal(mqttMsg)
	require.NoError(t, err)

	msg := &paho.Publish{
		Topic:   "cs/ocpp2.0.1/cs001",
		Payload: mqttMsgBytes,
	}

	handler(msg)

	assert.Equal(t, "cs001", chargeStationId)
	assert.Equal(t, mqtt.MessageTypeCallError, message.MessageType)
	assert.Equal(t, "Heartbeat", message.Action)
	assert.Equal(t, "1234", message.MessageId)
	assert.Equal(t, mqtt.ErrorFormatViolation, message.ErrorCode)
	assert.NotEmpty(t, message.ErrorDescription)
}
