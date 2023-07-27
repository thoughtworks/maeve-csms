// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"io/fs"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
)

type Handler struct {
	mqttBrokerUrls        []*url.URL
	mqttPrefix            string
	mqttGroup             string
	mqttConnectTimeout    time.Duration
	mqttConnectRetryDelay time.Duration
	mqttKeepAliveInterval uint16
	clock                 clock.PassiveClock
	tariffService         services.TariffService
	certValidationService services.CertificateValidationService
	certSignerService     services.CertificateSignerService
	certProviderService   services.EvCertificateProvider
	heartbeatInterval     time.Duration
	schemaFS              fs.FS
	storageEngine         store.Engine
	tracer                trace.Tracer
}

type HandlerOpt func(h *Handler)

func WithMqttBrokerUrl(brokerUrl *url.URL) HandlerOpt {
	return func(h *Handler) {
		h.mqttBrokerUrls = append(h.mqttBrokerUrls, brokerUrl)
	}
}

func WithMqttPrefix(mqttPrefix string) HandlerOpt {
	return func(h *Handler) {
		h.mqttPrefix = mqttPrefix
	}
}

func WithMqttGroup(mqttGroup string) HandlerOpt {
	return func(h *Handler) {
		h.mqttGroup = mqttGroup
	}
}

func WithMqttConnectSettings(mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval time.Duration) HandlerOpt {
	return func(handler *Handler) {
		handler.mqttConnectTimeout = mqttConnectTimeout
		handler.mqttConnectRetryDelay = mqttConnectRetryDelay
		handler.mqttKeepAliveInterval = uint16(mqttKeepAliveInterval.Round(time.Second).Seconds())
	}
}

func WithClock(clock clock.PassiveClock) HandlerOpt {
	return func(h *Handler) {
		h.clock = clock
	}
}

func WithTariffService(tariffService services.TariffService) HandlerOpt {
	return func(h *Handler) {
		h.tariffService = tariffService
	}
}

func WithCertValidationService(certValidationService services.CertificateValidationService) HandlerOpt {
	return func(h *Handler) {
		h.certValidationService = certValidationService
	}
}

func WithCertSignerService(certSignerService services.CertificateSignerService) HandlerOpt {
	return func(h *Handler) {
		h.certSignerService = certSignerService
	}
}

func WithCertificateProviderService(certProviderService services.EvCertificateProvider) HandlerOpt {
	return func(h *Handler) {
		h.certProviderService = certProviderService
	}
}

func WithHeartbeatInterval(heartbeatInterval time.Duration) HandlerOpt {
	return func(h *Handler) {
		h.heartbeatInterval = heartbeatInterval
	}
}

func WithSchemaFS(fs fs.FS) HandlerOpt {
	return func(h *Handler) {
		h.schemaFS = fs
	}
}

func WithStorageEngine(store store.Engine) HandlerOpt {
	return func(h *Handler) {
		h.storageEngine = store
	}
}

func WithOtelTracer(tracer trace.Tracer) HandlerOpt {
	return func(h *Handler) {
		h.tracer = tracer
	}
}

func NewHandler(opts ...HandlerOpt) *Handler {
	h := new(Handler)
	for _, opt := range opts {
		opt(h)
	}
	ensureDefaults(h)
	return h
}

func ensureDefaults(h *Handler) {
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
	if h.clock == nil {
		h.clock = clock.RealClock{}
	}
	if h.certValidationService == nil {
		h.certValidationService = services.OnlineCertificateValidationService{}
	}
	if h.heartbeatInterval == 0 {
		h.heartbeatInterval = time.Minute
	}
	if h.schemaFS == nil {
		h.schemaFS = schemas.OcppSchemas
	}
	if h.tracer == nil {
		h.tracer = trace.NewNoopTracerProvider().Tracer("")
	}
}

func (h *Handler) Connect(errCh chan error) {
	ctx, cancel := context.WithTimeout(context.Background(), h.mqttConnectTimeout)
	defer cancel()

	// ProxyEmitter supports late binding of the E implementation
	v16Emitter := &ProxyEmitter{}
	v201Emitter := &ProxyEmitter{}

	v16Router := NewV16Router(v16Emitter, h.clock, h.storageEngine, h.storageEngine, h.certValidationService, h.certSignerService, h.certProviderService, h.heartbeatInterval, h.schemaFS)
	v201Router := NewV201Router(v201Emitter, h.clock, h.storageEngine, h.storageEngine, h.tariffService, h.certValidationService, h.certSignerService, h.certProviderService, h.heartbeatInterval)

	mqttV16Topic := fmt.Sprintf("$share/%s/%s/in/ocpp1.6/#", h.mqttGroup, h.mqttPrefix)
	mqttV201Topic := fmt.Sprintf("$share/%s/%s/in/ocpp2.0.1/#", h.mqttGroup, h.mqttPrefix)

	readyCh := make(chan struct{})

	clientId := fmt.Sprintf("%s-%s", h.mqttGroup, randSeq(5))

	mqttRouter := paho.NewStandardRouter()
	mqttRouter.RegisterHandler(mqttV16Topic, NewGatewayMessageHandler(context.Background(), clientId, h.tracer, v16Router, v16Emitter, h.schemaFS))
	mqttRouter.RegisterHandler(mqttV201Topic, NewGatewayMessageHandler(context.Background(), clientId, h.tracer, v201Router, v201Emitter, h.schemaFS))

	var mqttConn *autopaho.ConnectionManager
	mqttConn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        h.mqttBrokerUrls,
		KeepAlive:         h.mqttKeepAliveInterval,
		ConnectRetryDelay: h.mqttConnectRetryDelay,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			v16Emitter.emitter = NewMqttEmitter(mqttConn, h.mqttPrefix, "ocpp1.6", h.tracer)
			v201Emitter.emitter = NewMqttEmitter(mqttConn, h.mqttPrefix, "ocpp2.0.1", h.tracer)

			_, err := mqttConn.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					mqttV16Topic:  {},
					mqttV201Topic: {},
				},
			})
			if err != nil {
				errCh <- err
			}

			slog.Info("subscribed to gateway+ocpp1.6", slog.String("brokerUrls", fmt.Sprintf("%+v", h.mqttBrokerUrls)), slog.String("topic", mqttV16Topic))
			slog.Info("subscribed to gateway+ocpp2.0.1", slog.String("brokerUrls", fmt.Sprintf("%+v", h.mqttBrokerUrls)), slog.String("topic", mqttV201Topic))

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

func NewGatewayMessageHandler(
	ctx context.Context,
	clientId string,
	tracer trace.Tracer,
	router *Router,
	emitter Emitter,
	schemaFS fs.FS) func(mqttMsg *paho.Publish) {
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
		chargeStationId := topicParts[len(topicParts)-1]

		var msg Message
		err := json.Unmarshal(mqttMsg.Payload, &msg)
		if err != nil {
			errMsg := NewErrorMessage("", "-1", ErrorInternalError, err)
			err = emitter.Emit(newCtx, chargeStationId, errMsg)
			if err != nil {
				slog.Error("unable to emit error message", "err", err)
			}
		}

		span.SetAttributes(
			attribute.String("csId", chargeStationId),
			attribute.String(fmt.Sprintf("%s.action", msg.MessageType), msg.Action),
			semconv.MessagingMessageConversationID(msg.MessageId),
		)

		if msg.MessageType == MessageTypeCallError {
			span.SetAttributes(
				attribute.String(fmt.Sprintf("%s.code", msg.MessageType), string(msg.ErrorCode)),
				attribute.String(fmt.Sprintf("%s.description", msg.MessageType), msg.ErrorDescription))
		}

		err = router.Route(newCtx, chargeStationId, msg, emitter, schemaFS)
		if err != nil {
			slog.Error("unable to route message", slog.String("chargeStationId", chargeStationId), slog.String("action", msg.Action), "err", err)
			span.SetStatus(codes.Error, "routing request failed")
			span.RecordError(err)
			var mqttError *Error
			var errMsg *Message
			if errors.As(err, &mqttError) {
				errMsg = NewErrorMessage(msg.Action, msg.MessageId, mqttError.ErrorCode, mqttError.wrappedError)
			} else {
				errMsg = NewErrorMessage(msg.Action, msg.MessageId, ErrorInternalError, err)
			}
			err = emitter.Emit(ctx, chargeStationId, errMsg)
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
