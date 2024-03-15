package mqtt_test

import (
	"context"
	"encoding/json"
	server "github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/testutil"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"github.com/thoughtworks/maeve-csms/manager/transport/mqtt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"strings"
	"testing"
	"time"
)

func TestListenerProcessesMessagesReceivedFromTheBroker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start the broker
	broker, clientUrl := mqtt.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()
	err := broker.Serve()
	require.NoError(t, err)

	// setup the handler
	receivedMsgCh := make(chan struct{})
	handler := func(ctx context.Context, chargeStationId string, msg *transport.Message) {
		assert.Equal(t, "cs001", chargeStationId)
		assert.Equal(t, transport.MessageTypeCall, msg.MessageType)
		assert.Equal(t, "Test", msg.Action)
		assert.Equal(t, "my-message-id", msg.MessageId)
		assert.Equal(t, json.RawMessage(`{"someKey":"someValue"}`), msg.RequestPayload)
		receivedMsgCh <- struct{}{}
	}

	// connect the listener to the broker
	listener := mqtt.NewListener(mqtt.WithMqttBrokerUrl[mqtt.Listener](clientUrl))
	conn, err := listener.Connect(ctx, transport.OcppVersion201, nil, transport.MessageHandlerFunc(handler))
	require.NoError(t, err)
	defer func() {
		if conn != nil {
			err := conn.Disconnect(ctx)
			require.NoError(t, err)
		}
	}()

	// publish message
	publishMessage(t, ctx, broker, transport.Message{
		MessageType:    transport.MessageTypeCall,
		Action:         "Test",
		MessageId:      "my-message-id",
		RequestPayload: json.RawMessage(`{"someKey":"someValue"}`),
	})

	// wait for message to be received / timeout
	select {
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test to complete")
	case <-receivedMsgCh:
		// do nothing
	}
}

func TestListenerExtractsTraceIdFromPayload(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tracer, exporter := testutil.GetTracer()

	// start the broker
	broker, clientUrl := mqtt.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()
	err := broker.Serve()
	require.NoError(t, err)

	// setup the handler
	receivedMsgCh := make(chan struct{})
	handler := func(ctx context.Context, chargeStationId string, msg *transport.Message) {
		receivedMsgCh <- struct{}{}
	}

	// connect the listener to the broker
	listener := mqtt.NewListener(
		mqtt.WithMqttBrokerUrl[mqtt.Listener](clientUrl),
		mqtt.WithOtelTracer[mqtt.Listener](tracer))
	conn, err := listener.Connect(ctx, transport.OcppVersion201, nil, transport.MessageHandlerFunc(handler))
	require.NoError(t, err)
	defer func() {
		if conn != nil {
			err := conn.Disconnect(ctx)
			require.NoError(t, err)
		}
	}()

	// publish message
	newCtx, span := tracer.Start(ctx, "test span")
	defer span.End()
	publishMessage(t, newCtx, broker, transport.Message{
		MessageType:    transport.MessageTypeCall,
		Action:         "Test",
		MessageId:      "my-message-id",
		RequestPayload: json.RawMessage(`{"someKey":"someValue"}`),
	})

	// wait for message to be received / timeout
	select {
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test to complete")
	case <-receivedMsgCh:
		require.Greater(t, len(exporter.GetSpans()), 0)
		assert.True(t, exporter.GetSpans()[0].Parent.HasTraceID())
	}
}

func TestListenerAddsTraceInformation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tracer, exporter := testutil.GetTracer()

	// start the broker
	broker, clientUrl := mqtt.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()
	err := broker.Serve()
	require.NoError(t, err)

	// setup the handler
	receivedMsgCh := make(chan struct{})
	handler := func(ctx context.Context, chargeStationId string, msg *transport.Message) {
		receivedMsgCh <- struct{}{}
	}

	// connect the listener to the broker
	listener := mqtt.NewListener(
		mqtt.WithMqttBrokerUrl[mqtt.Listener](clientUrl),
		mqtt.WithOtelTracer[mqtt.Listener](tracer))
	conn, err := listener.Connect(ctx, transport.OcppVersion201, nil, transport.MessageHandlerFunc(handler))
	require.NoError(t, err)
	defer func() {
		if conn != nil {
			err := conn.Disconnect(ctx)
			require.NoError(t, err)
		}
	}()

	// publish message
	newCtx, span := tracer.Start(ctx, "test span")
	defer span.End()
	publishMessage(t, newCtx, broker, transport.Message{
		MessageType:    transport.MessageTypeCall,
		Action:         "Test",
		MessageId:      "my-message-id",
		RequestPayload: json.RawMessage(`{"someKey":"someValue"}`),
	})

	// wait for message to be received / timeout
	select {
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test to complete")
	case <-receivedMsgCh:
		require.Greater(t, len(exporter.GetSpans()), 0)
		testutil.AssertSpan(t, &exporter.GetSpans()[0], "cs/in/ocpp2.0.1/# receive", map[string]any{
			"messaging.system":                     "mqtt",
			"messaging.operation":                  "receive",
			"messaging.message.payload_size_bytes": 81,
			"ocpp.version":                         "2.0.1",
			"csId":                                 "cs001",
			"call.action":                          "Test",
			"messaging.consumer.id": func(val attribute.Value) bool {
				return strings.HasPrefix(val.AsString(), "manager-")
			},
			"messaging.message.conversation_id": func(val attribute.Value) bool {
				return val.AsString() != ""
			},
		})
	}
}

func publishMessage(t *testing.T, ctx context.Context, broker *server.Server, msg transport.Message) {
	msgBytes, err := json.Marshal(msg)
	require.NoError(t, err)

	correlationMap := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(correlationMap))
	correlationData, err := json.Marshal(correlationMap)
	require.NoError(t, err)

	cl := broker.NewClient(nil, "local", "inline", true)
	err = broker.InjectPacket(cl, packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type:   packets.Publish,
			Qos:    0,
			Retain: false,
		},
		TopicName: "cs/in/ocpp2.0.1/cs001",
		Payload:   msgBytes,
		PacketID:  uint16(0),
		Properties: packets.Properties{
			CorrelationData: correlationData,
		},
	})
	require.NoError(t, err)
}
