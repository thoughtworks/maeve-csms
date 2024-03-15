// SPDX-License-Identifier: Apache-2.0

package mqtt_test

import (
	"context"
	"encoding/json"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	mqtt2 "github.com/thoughtworks/maeve-csms/manager/transport/mqtt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"net/url"
	"testing"
	"time"
)

func TestEmitterSendsOcpp201Message(t *testing.T) {
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

	// subscribe to the output channel
	rcvdCh := make(chan struct{})

	mqttClient := listenForMessageSentByManager(t, ctx, clientUrl, paho.NewSingleHandlerRouter(func(publish *paho.Publish) {
		assert.Equal(t, "cs/out/ocpp2.0.1/cs001", publish.Topic)
		var msg transport.Message
		err := json.Unmarshal(publish.Payload, &msg)
		assert.NoError(t, err, "payload is not the expected message type")
		assert.Equal(t, transport.MessageTypeCall, msg.MessageType)
		assert.Equal(t, "TriggerMessage", msg.Action)
		assert.Equal(t, "1234", msg.MessageId)
		assert.JSONEq(t, `{"requestedMessage":"Heartbeat"}`, string(msg.RequestPayload))
		rcvdCh <- struct{}{}
	}))

	defer func() {
		_ = mqttClient.Disconnect(ctx)
	}()

	// publish a message to the input channel
	msg := transport.Message{
		MessageType:    transport.MessageTypeCall,
		Action:         "TriggerMessage",
		MessageId:      "1234",
		RequestPayload: []byte(`{"requestedMessage":"Heartbeat"}`),
	}
	err = emitter.Emit(context.Background(), transport.OcppVersion201, "cs001", &msg)
	require.NoError(t, err)

	// wait for success
	select {
	case <-rcvdCh:
		// success
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test")
	}
}

func TestEmitterSendsOcpp16Message(t *testing.T) {
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

	// subscribe to the output channel
	rcvdCh := make(chan struct{})

	mqttClient := listenForMessageSentByManager(t, ctx, clientUrl, paho.NewSingleHandlerRouter(func(publish *paho.Publish) {
		assert.Equal(t, "cs/out/ocpp1.6/cs001", publish.Topic)
		var msg transport.Message
		err := json.Unmarshal(publish.Payload, &msg)
		assert.NoError(t, err, "payload is not the expected message type")
		assert.Equal(t, transport.MessageTypeCall, msg.MessageType)
		assert.Equal(t, "TriggerMessage", msg.Action)
		assert.Equal(t, "1234", msg.MessageId)
		assert.JSONEq(t, `{"requestedMessage":"Heartbeat"}`, string(msg.RequestPayload))
		rcvdCh <- struct{}{}
	}))

	defer func() {
		_ = mqttClient.Disconnect(ctx)
	}()

	// publish a message to the input channel
	msg := transport.Message{
		MessageType:    transport.MessageTypeCall,
		Action:         "TriggerMessage",
		MessageId:      "1234",
		RequestPayload: []byte(`{"requestedMessage":"Heartbeat"}`),
	}
	err = emitter.Emit(context.Background(), transport.OcppVersion16, "cs001", &msg)
	require.NoError(t, err)

	// wait for success
	select {
	case <-rcvdCh:
		// success
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test")
	}
}

func TestEmitterAddsCorrelationData(t *testing.T) {
	traceExporter := tracetest.NewInMemoryExporter()
	tracerProvider := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithSyncer(traceExporter),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tracer := tracerProvider.Tracer("test")

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
		mqtt2.WithMqttPrefix[mqtt2.Emitter]("cs"),
		mqtt2.WithOtelTracer[mqtt2.Emitter](tracer))

	// subscribe to the output channel
	rcvdCh := make(chan struct{})

	mqttClient := listenForMessageSentByManager(t, ctx, clientUrl, paho.NewSingleHandlerRouter(func(publish *paho.Publish) {
		assert.Equal(t, "cs/out/ocpp2.0.1/cs001", publish.Topic)
		require.NotNil(t, publish.Properties)
		require.NotNil(t, publish.Properties.CorrelationData)
		correlationMap := make(map[string]string)
		err := json.Unmarshal(publish.Properties.CorrelationData, &correlationMap)
		require.NoError(t, err)
		assert.NotEmpty(t, correlationMap["traceparent"])

		rcvdCh <- struct{}{}
	}))

	defer func() {
		_ = mqttClient.Disconnect(ctx)
	}()

	// publish a message to the input channel
	msg := transport.Message{
		MessageType:    transport.MessageTypeCall,
		Action:         "TriggerMessage",
		MessageId:      "1234",
		RequestPayload: []byte(`{"requestedMessage":"Heartbeat"}`),
	}

	err = emitter.Emit(ctx, transport.OcppVersion201, "cs001", &msg)
	require.NoError(t, err)

	// wait for success
	select {
	case <-rcvdCh:
		// success
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test")
	}
}

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
