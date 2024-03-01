// SPDX-License-Identifier: Apache-2.0

package mqtt_test

import (
	"context"
	"encoding/json"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	server "github.com/mochi-co/mqtt/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	mqtt2 "github.com/thoughtworks/maeve-csms/manager/transport/mqtt"
	"k8s.io/utils/clock"
	"net/url"
	"testing"
	"time"
)

func listenForMessageSentByManager(t *testing.T, ctx context.Context, clientUrl *url.URL, router paho.Router) *autopaho.ConnectionManager {
	mqttClient, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{clientUrl},
		KeepAlive:         10,
		ConnectRetryDelay: 10,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err := manager.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					"cs/out/ocpp1.6/cs001":   {},
					"cs/out/ocpp2.0.1/cs001": {},
				},
			})
			require.NoError(t, err)
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "test",
			Router:   router,
		},
	})
	require.NoError(t, err)

	err = mqttClient.AwaitConnection(ctx)
	require.NoError(t, err)

	return mqttClient
}

func publishMessageToManager(t *testing.T, broker *server.Server, msg transport.Message) {
	msgBytes, err := json.Marshal(msg)
	require.NoError(t, err)

	err = broker.Publish("cs/in/ocpp2.0.1/cs001", msgBytes, false, 0)
	require.NoError(t, err)
}

func TestMqttReceiverRespondsToACall(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start the broker
	broker, clientUrl := mqtt2.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()

	err = broker.Serve()
	require.NoError(t, err)

	emitter := mqtt2.NewEmitter(
		mqtt2.WithMqttBrokerUrl[mqtt2.Emitter](clientUrl),
		mqtt2.WithMqttPrefix[mqtt2.Emitter]("cs"))

	router := &handlers.Router{
		Emitter:     emitter,
		SchemaFS:    schemas.OcppSchemas,
		OcppVersion: transport.OcppVersion201,
		CallRoutes: map[string]handlers.CallRoute{
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				RequestSchema:  "ocpp201/HeartbeatRequest.json",
				ResponseSchema: "ocpp201/HeartbeatResponse.json",
				Handler: handlers.CallHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error) {
					return &types.HeartbeatResponseJson{
						CurrentTime: now.Format(time.RFC3339),
					}, nil
				}),
			},
		},
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}

	// connect the mqtt handler
	handler := mqtt2.NewReceiver(
		mqtt2.WithMqttBrokerUrl[mqtt2.Receiver](clientUrl),
		mqtt2.WithMqttPrefix[mqtt2.Receiver]("cs"),
		mqtt2.WithRouter(router),
		mqtt2.WithEmitter(emitter))
	errCh := make(chan error)
	handler.Connect(errCh)

	// subscribe to the output channel
	rcvdCh := make(chan struct{})

	mqttClient := listenForMessageSentByManager(t, ctx, clientUrl, paho.NewSingleHandlerRouter(func(publish *paho.Publish) {
		assert.Equal(t, "cs/out/ocpp2.0.1/cs001", publish.Topic)
		var msg transport.Message
		err := json.Unmarshal(publish.Payload, &msg)
		assert.NoError(t, err, "payload is not the expected message type")
		assert.Equal(t, transport.MessageTypeCallResult, msg.MessageType)
		assert.Equal(t, "Heartbeat", msg.Action)
		assert.Equal(t, "1234", msg.MessageId)
		assert.JSONEq(t, `{"currentTime":"2023-06-15T15:05:00+01:00"}`, string(msg.ResponsePayload))
		rcvdCh <- struct{}{}
	}))

	defer func() {
		_ = mqttClient.Disconnect(ctx)
	}()

	// publish a message to the input channel
	msg := transport.Message{
		MessageType:    transport.MessageTypeCall,
		Action:         "Heartbeat",
		MessageId:      "1234",
		RequestPayload: []byte("{}"),
	}
	publishMessageToManager(t, broker, msg)

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

func TestMqttReceiverEmitsErrorWhenFailingToUnmarshallIncomingMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start the broker
	broker, clientUrl := mqtt2.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()

	err := broker.Serve()
	require.NoError(t, err)

	emitter := mqtt2.NewEmitter(
		mqtt2.WithMqttBrokerUrl[mqtt2.Emitter](clientUrl),
		mqtt2.WithMqttPrefix[mqtt2.Emitter]("cs"))

	router := &handlers.Router{
		Emitter:          emitter,
		SchemaFS:         schemas.OcppSchemas,
		OcppVersion:      transport.OcppVersion201,
		CallRoutes:       map[string]handlers.CallRoute{},
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}

	// connect the mqtt handler
	handler := mqtt2.NewReceiver(
		mqtt2.WithMqttBrokerUrl[mqtt2.Receiver](clientUrl),
		mqtt2.WithMqttPrefix[mqtt2.Receiver]("cs"),
		mqtt2.WithRouter(router),
		mqtt2.WithEmitter(emitter))
	errCh := make(chan error)
	handler.Connect(errCh)

	// subscribe to the output channel
	rcvdCh := make(chan struct{})

	mqttClient := listenForMessageSentByManager(t, ctx, clientUrl, paho.NewSingleHandlerRouter(func(publish *paho.Publish) {
		assert.Equal(t, "cs/out/ocpp2.0.1/cs001", publish.Topic)
		var msg transport.Message
		err := json.Unmarshal(publish.Payload, &msg)
		assert.NoError(t, err, "payload is not the expected message type")
		assert.Equal(t, transport.MessageTypeCallError, msg.MessageType)
		assert.Equal(t, "", msg.Action)
		assert.Equal(t, "-1", msg.MessageId)
		assert.Equal(t, transport.ErrorInternalError, msg.ErrorCode)
		rcvdCh <- struct{}{}
	}))

	defer func() {
		_ = mqttClient.Disconnect(ctx)
	}()

	// publish the message to the handler
	err = broker.Publish("cs/in/ocpp2.0.1/cs001",
		[]byte(""), false, 0)
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

func TestMqttReceiverEmitsErrorWithCodeFromRouter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start the broker
	broker, clientUrl := mqtt2.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()

	err := broker.Serve()
	require.NoError(t, err)

	emitter := mqtt2.NewEmitter(
		mqtt2.WithMqttBrokerUrl[mqtt2.Emitter](clientUrl),
		mqtt2.WithMqttPrefix[mqtt2.Emitter]("cs"))

	router := &handlers.Router{
		Emitter:          emitter,
		SchemaFS:         schemas.OcppSchemas,
		OcppVersion:      transport.OcppVersion201,
		CallRoutes:       map[string]handlers.CallRoute{},
		CallResultRoutes: map[string]handlers.CallResultRoute{},
	}

	// connect the mqtt handler
	handler := mqtt2.NewReceiver(
		mqtt2.WithMqttBrokerUrl[mqtt2.Receiver](clientUrl),
		mqtt2.WithMqttPrefix[mqtt2.Receiver]("cs"),
		mqtt2.WithRouter(router),
		mqtt2.WithEmitter(emitter))
	errCh := make(chan error)
	handler.Connect(errCh)

	// subscribe to the output channel
	rcvdCh := make(chan struct{})

	mqttClient := listenForMessageSentByManager(t, ctx, clientUrl, paho.NewSingleHandlerRouter(func(publish *paho.Publish) {
		assert.Equal(t, "cs/out/ocpp2.0.1/cs001", publish.Topic)
		var msg transport.Message
		err := json.Unmarshal(publish.Payload, &msg)
		assert.NoError(t, err, "payload is not the expected message type")
		assert.Equal(t, transport.MessageTypeCallError, msg.MessageType)
		assert.Equal(t, "Heartbeat", msg.Action)
		assert.Equal(t, "1234", msg.MessageId)
		assert.Equal(t, transport.ErrorNotImplemented, msg.ErrorCode)
		rcvdCh <- struct{}{}
	}))

	defer func() {
		_ = mqttClient.Disconnect(ctx)
	}()

	// publish a message to the input channel
	msg := transport.Message{
		MessageType:    transport.MessageTypeCall,
		Action:         "Heartbeat",
		MessageId:      "1234",
		RequestPayload: []byte("{}"),
	}
	publishMessageToManager(t, broker, msg)

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

func TestMqttReceiverEmitsErrorFromRouter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start the broker
	broker, clientUrl := mqtt2.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()

	err := broker.Serve()
	require.NoError(t, err)

	emitter := mqtt2.NewEmitter(
		mqtt2.WithMqttBrokerUrl[mqtt2.Emitter](clientUrl),
		mqtt2.WithMqttPrefix[mqtt2.Emitter]("cs"))

	router := &handlers.Router{
		Emitter:     emitter,
		SchemaFS:    schemas.OcppSchemas,
		OcppVersion: transport.OcppVersion201,
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

	// connect the mqtt handler
	handler := mqtt2.NewReceiver(
		mqtt2.WithMqttBrokerUrl[mqtt2.Receiver](clientUrl),
		mqtt2.WithMqttPrefix[mqtt2.Receiver]("cs"),
		mqtt2.WithRouter(router),
		mqtt2.WithEmitter(emitter))
	errCh := make(chan error)
	handler.Connect(errCh)

	// subscribe to the output channel
	rcvdCh := make(chan struct{})

	mqttClient := listenForMessageSentByManager(t, ctx, clientUrl, paho.NewSingleHandlerRouter(func(publish *paho.Publish) {
		assert.Equal(t, "cs/out/ocpp2.0.1/cs001", publish.Topic)
		var msg transport.Message
		err := json.Unmarshal(publish.Payload, &msg)
		assert.NoError(t, err, "payload is not the expected message type")
		assert.Equal(t, transport.MessageTypeCallError, msg.MessageType)
		assert.Equal(t, "Heartbeat", msg.Action)
		assert.Equal(t, "1234", msg.MessageId)
		assert.Equal(t, transport.ErrorFormatViolation, msg.ErrorCode)
		rcvdCh <- struct{}{}
	}))

	defer func() {
		_ = mqttClient.Disconnect(ctx)
	}()

	// publish a message to the input channel
	msg := transport.Message{
		MessageType:    transport.MessageTypeCall,
		Action:         "Heartbeat",
		MessageId:      "1234",
		RequestPayload: []byte("{}"),
	}
	publishMessageToManager(t, broker, msg)

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
