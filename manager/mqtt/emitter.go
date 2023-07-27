// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type Emitter interface {
	Emit(ctx context.Context, chargeStationId string, message *Message) error
}

type EmitterFunc func(ctx context.Context, chargeStationId string, message *Message) error

func (e EmitterFunc) Emit(ctx context.Context, chargeStationId string, message *Message) error {
	return e(ctx, chargeStationId, message)
}

type ProxyEmitter struct {
	emitter Emitter
}

func (p *ProxyEmitter) Emit(ctx context.Context, chargeStationId string, message *Message) error {
	if p.emitter == nil {
		return fmt.Errorf("no emitter configured")
	}
	return p.emitter.Emit(ctx, chargeStationId, message)
}

type MqttEmitter struct {
	conn        *autopaho.ConnectionManager
	mqttPrefix  string
	ocppVersion string
	tracer      trace.Tracer
}

func (m *MqttEmitter) Emit(ctx context.Context, chargeStationId string, message *Message) error {
	topic := fmt.Sprintf("%s/out/%s/%s", m.mqttPrefix, m.ocppVersion, chargeStationId)
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshalling response of type %s: %v", message.Action, err)
	}

	newCtx, span := m.tracer.Start(ctx,
		fmt.Sprintf("%s/out/%s/# publish", m.mqttPrefix, m.ocppVersion),
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

	_, err = m.conn.Publish(newCtx, &paho.Publish{
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

func NewMqttEmitter(conn *autopaho.ConnectionManager, mqttPrefix, ocppVersion string, tracer trace.Tracer) Emitter {
	return &MqttEmitter{
		conn:        conn,
		mqttPrefix:  mqttPrefix,
		ocppVersion: ocppVersion,
		tracer:      tracer,
	}
}
