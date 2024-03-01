// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"golang.org/x/exp/slog"
)

// Receiver is an implementation of transport.Receiver that uses MQTT
// as the transport.
//
// The Receiver subscribes to topics matching that pattern:
// <prefix>/in/<ocpp-version>/#. The prefix is configured, or defaults
// to `cs`. There will be a separate subscription for each configured
// router - which defines the ocpp-version.
//
// The Receiver subscribes using an MQTT 5 group subscription. The
// group name can be configured, but defaults to `manager`.
//
// The Receiver defaults to connecting to the broker on 127.0.0.1:1883.
type Receiver struct {
	connection
	mqttGroup string
	tracer    trace.Tracer
	routers   map[transport.OcppVersion]transport.Router
	emitter   transport.Emitter
}

func NewReceiver(opts ...Opt[Receiver]) *Receiver {
	h := new(Receiver)
	h.routers = make(map[transport.OcppVersion]transport.Router)
	for _, opt := range opts {
		opt(h)
	}
	ensureHandlerDefaults(h)
	return h
}

func ensureHandlerDefaults(h *Receiver) {
	if h.mqttBrokerUrls == nil {
		u, err := url.Parse("mqtt://127.0.0.1:1883/")
		if err != nil {
			panic(err)
		}
		h.mqttBrokerUrls = []*url.URL{u}
	}
	if h.mqttPrefix == "" {
		h.mqttPrefix = "cs"
	}
	if h.mqttGroup == "" {
		h.mqttGroup = "manager"
	}
	if h.mqttConnectTimeout == 0 {
		h.mqttConnectTimeout = 10 * time.Second
	}
	if h.mqttConnectRetryDelay == 0 {
		h.mqttConnectRetryDelay = 1 * time.Second
	}
	if h.mqttKeepAliveInterval == 0 {
		h.mqttKeepAliveInterval = 10
	}
	if h.tracer == nil {
		h.tracer = noop.NewTracerProvider().Tracer("")
	}
	if h.emitter == nil {
		h.emitter = NewEmitter(
			WithMqttBrokerUrls[Emitter](h.mqttBrokerUrls),
			WithMqttPrefix[Emitter](h.mqttPrefix),
			WithMqttConnectSettings[Emitter](h.mqttConnectTimeout, h.mqttConnectRetryDelay, time.Duration(h.mqttKeepAliveInterval)*time.Second),
			WithOtelTracer[Emitter](h.tracer))
	}
}

func (h *Receiver) Connect(errCh chan error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.mqttConnectTimeout)
	defer cancel()

	readyCh := make(chan struct{})

	clientId := fmt.Sprintf("%s-%s", h.mqttGroup, randSeq(5))

	mqttRouter := paho.NewStandardRouter()
	var mqttConn *autopaho.ConnectionManager
	mqttConn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        h.mqttBrokerUrls,
		KeepAlive:         h.mqttKeepAliveInterval,
		ConnectRetryDelay: h.mqttConnectRetryDelay,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			subscriptions := make(map[string]paho.SubscribeOptions)

			for version, router := range h.routers {
				topic := fmt.Sprintf("$share/%s/%s/in/%s/#", h.mqttGroup, h.mqttPrefix, version)
				subscriptions[topic] = paho.SubscribeOptions{}
				mqttRouter.RegisterHandler(topic, newGatewayMessageHandler(
					context.Background(),
					clientId,
					h.tracer,
					version,
					router,
					h.emitter))
				slog.Info("subscribed to gateway", slog.String("brokerUrls", fmt.Sprintf("%+v", h.mqttBrokerUrls)), slog.String("topic", topic))
			}

			_, err := mqttConn.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: subscriptions,
			})
			if err != nil {
				errCh <- err
			}

			readyCh <- struct{}{}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: clientId,
			Router:   mqttRouter,
		},
	})
	if err != nil {
		errCh <- err
		return
	}

	select {
	case <-ctx.Done():
		errCh <- errors.New("timeout waiting for mqtt connection setup")
	case <-readyCh:
		// do nothing
	}
}

func getTopicPattern(topic string) string {
	parts := strings.Split(topic, "/")
	parts[len(parts)-1] = "#"
	return strings.Join(parts, "/")
}

func newGatewayMessageHandler(
	ctx context.Context,
	clientId string,
	tracer trace.Tracer,
	ocppVersion transport.OcppVersion,
	router transport.Router,
	emitter transport.Emitter) func(mqttMsg *paho.Publish) {
	return func(mqttMsg *paho.Publish) {
		baseCtx := context.Background()

		if mqttMsg.Properties != nil && mqttMsg.Properties.CorrelationData != nil {
			correlationMap := make(map[string]string)
			err := json.Unmarshal(mqttMsg.Properties.CorrelationData, &correlationMap)
			if err != nil {
				slog.Warn("failed to unmarshal correlation data", "error", err)
			} else {
				baseCtx = otel.GetTextMapPropagator().Extract(baseCtx, propagation.MapCarrier(correlationMap))
			}
		}

		newCtx, span := tracer.Start(baseCtx,
			fmt.Sprintf("%s receive", getTopicPattern(mqttMsg.Topic)),
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				semconv.MessagingSystem("mqtt"),
				semconv.MessagingConsumerID(clientId),
				semconv.MessagingMessagePayloadSizeBytes(len(mqttMsg.Payload)),
				semconv.MessagingOperationKey.String("receive"),
			))
		defer span.End()

		topicParts := strings.Split(mqttMsg.Topic, "/")
		var chargeStationId = topicParts[len(topicParts)-1]
		var msg transport.Message
		err := json.Unmarshal(mqttMsg.Payload, &msg)
		if err != nil {
			errMsg := transport.NewErrorMessage("", "-1", transport.ErrorInternalError, err)
			err = emitter.Emit(newCtx, ocppVersion, chargeStationId, errMsg)
			if err != nil {
				slog.Error("unable to emit error message", "err", err)
			}
		}

		span.SetAttributes(
			attribute.String("csId", chargeStationId),
			attribute.String(fmt.Sprintf("%s.action", msg.MessageType), msg.Action),
			semconv.MessagingMessageConversationID(msg.MessageId),
		)

		if msg.MessageType == transport.MessageTypeCallError {
			span.SetAttributes(
				attribute.String(fmt.Sprintf("%s.code", msg.MessageType), string(msg.ErrorCode)),
				attribute.String(fmt.Sprintf("%s.description", msg.MessageType), msg.ErrorDescription))
		}

		err = router.Route(newCtx, chargeStationId, msg)
		if err != nil {
			slog.Error("unable to route message", slog.String("chargeStationId", chargeStationId), slog.String("action", msg.Action), "err", err)
			span.SetStatus(codes.Error, "routing request failed")
			span.RecordError(err)
			var mqttError *transport.Error
			var errMsg *transport.Message
			if errors.As(err, &mqttError) {
				errMsg = transport.NewErrorMessage(msg.Action, msg.MessageId, mqttError.ErrorCode, mqttError.WrappedError)
			} else {
				errMsg = transport.NewErrorMessage(msg.Action, msg.MessageId, transport.ErrorInternalError, err)
			}
			err = emitter.Emit(ctx, ocppVersion, chargeStationId, errMsg)
			if err != nil {
				slog.Error("unable to emit error message", "err", err)
			}
		} else {
			span.SetStatus(codes.Ok, "ok")
		}
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		//#nosec G404 - client suffix does not require secure random number generator
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
