// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"net/url"
	"sync"
	"time"
)

// Emitter is an implementation of transport.Emitter that uses MQTT
// as the transport.
//
// Messages are published on a topic that is composed of a number of
// elements: <prefix>/out/<ocpp-version>/<cs-id>. The prefix is
// configured, the ocpp-version and cs-id are provided to the Emit
// function. If not configured the default prefix is `cs`.
//
// The Emitter defaults to connecting to a broker on 127.0.0.1:1883.
type Emitter struct {
	sync.Mutex
	connection
	tracer trace.Tracer
	conn   *autopaho.ConnectionManager
}

func NewEmitter(opts ...Opt[Emitter]) transport.Emitter {
	e := new(Emitter)
	for _, opt := range opts {
		opt(e)
	}
	ensureEmitterDefaults(e)
	return e
}

func (e *Emitter) Emit(ctx context.Context, ocppVersion transport.OcppVersion, chargeStationId string, message *transport.Message) error {
	topic := fmt.Sprintf("%s/out/%s/%s", e.mqttPrefix, ocppVersion, chargeStationId)
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshalling response of type %s: %v", message.Action, err)
	}

	newCtx, span := e.tracer.Start(ctx,
		fmt.Sprintf("%s/out/%s/# publish", e.mqttPrefix, ocppVersion),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.MessagingSystem("mqtt"),
			semconv.MessagingMessagePayloadSizeBytes(len(payload)),
			semconv.MessagingOperationKey.String("publish"),
			semconv.MessagingMessageConversationID(message.MessageId),
			attribute.String("csId", chargeStationId),
		))
	defer span.End()

	correlationMap := make(map[string]string)
	otel.GetTextMapPropagator().Inject(newCtx, propagation.MapCarrier(correlationMap))

	correlationData, err := json.Marshal(correlationMap)
	if err != nil {
		return fmt.Errorf("marshalling correlation map: %v", err)
	}

	err = e.ensureConnection(ctx)
	if err != nil {
		return fmt.Errorf("connecting to MQTT: %v", err)
	}

	_, err = e.conn.Publish(newCtx, &paho.Publish{
		Topic:   topic,
		Payload: payload,
		Properties: &paho.PublishProperties{
			CorrelationData: correlationData,
		},
	})
	if err != nil {
		return fmt.Errorf("publishing to %s: %v", topic, err)
	}
	return nil
}

func ensureEmitterDefaults(e *Emitter) {
	if e.mqttBrokerUrls == nil {
		u, err := url.Parse("mqtt://127.0.0.1:1883/")
		if err != nil {
			panic(err)
		}
		e.mqttBrokerUrls = []*url.URL{u}
	}
	if e.mqttPrefix == "" {
		e.mqttPrefix = "cs"
	}
	if e.mqttConnectTimeout == 0 {
		e.mqttConnectTimeout = 10 * time.Second
	}
	if e.mqttConnectRetryDelay == 0 {
		e.mqttConnectRetryDelay = 1 * time.Second
	}
	if e.mqttKeepAliveInterval == 0 {
		e.mqttKeepAliveInterval = 10
	}
	if e.tracer == nil {
		e.tracer = noop.NewTracerProvider().Tracer("")
	}
}

func (e *Emitter) ensureConnection(ctx context.Context) error {
	e.Lock()
	defer e.Unlock()
	if e.conn == nil {
		conn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
			BrokerUrls:        e.mqttBrokerUrls,
			KeepAlive:         e.mqttKeepAliveInterval,
			ConnectRetryDelay: e.mqttConnectRetryDelay,
			ClientConfig: paho.ClientConfig{
				ClientID: fmt.Sprintf("%s-%s", "manager-emit", randSeq(5)),
			},
		})
		if err != nil {
			return err
		}
		e.conn = conn

		err = conn.AwaitConnection(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
