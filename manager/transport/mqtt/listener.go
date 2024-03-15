package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"golang.org/x/exp/slog"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type Listener struct {
	connectionDetails
	mqttGroup string
	tracer    trace.Tracer
}

func NewListener(opts ...Opt[Listener]) *Listener {
	l := new(Listener)
	for _, opt := range opts {
		opt(l)
	}
	ensureListenerDefaults(l)
	return l
}

func ensureListenerDefaults(l *Listener) {
	if l.mqttBrokerUrls == nil {
		u, err := url.Parse("mqtt://127.0.0.1:1883/")
		if err != nil {
			panic(err)
		}
		l.mqttBrokerUrls = []*url.URL{u}
	}
	if l.mqttPrefix == "" {
		l.mqttPrefix = "cs"
	}
	if l.mqttGroup == "" {
		l.mqttGroup = "manager"
	}
	if l.mqttConnectTimeout == 0 {
		l.mqttConnectTimeout = 10 * time.Second
	}
	if l.mqttConnectRetryDelay == 0 {
		l.mqttConnectRetryDelay = 1 * time.Second
	}
	if l.mqttKeepAliveInterval == 0 {
		l.mqttKeepAliveInterval = 10
	}
	if l.tracer == nil {
		l.tracer = noop.NewTracerProvider().Tracer("")
	}
}

func (l *Listener) Connect(ctx context.Context, ocppVersion transport.OcppVersion, chargeStationId *string, handler transport.MessageHandler) (transport.Connection, error) {
	var err error

	ctx, cancel := context.WithTimeout(ctx, l.mqttConnectTimeout)
	defer cancel()

	clientId := fmt.Sprintf("%s-%s", l.mqttGroup, randSeq(5))

	readyCh := make(chan struct{})

	var topic string
	if chargeStationId != nil {
		topic = fmt.Sprintf("%s/in/%s/%s", l.mqttPrefix, ocppVersion, *chargeStationId)
	} else {
		topic = fmt.Sprintf("$share/%s/%s/in/%s/#", l.mqttGroup, l.mqttPrefix, ocppVersion)
	}

	conn := new(connection)
	mqttRouter := paho.NewStandardRouter()
	conn.mqttConn, err = autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        l.mqttBrokerUrls,
		KeepAlive:         l.mqttKeepAliveInterval,
		ConnectRetryDelay: l.mqttConnectRetryDelay,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err := manager.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					topic: {},
				},
			})
			if err != nil {
				slog.Error("failed to subscribe to topic", "topic", topic)
				return
			}
			mqttRouter.UnregisterHandler(topic)
			mqttRouter.RegisterHandler(topic, func(mqttMsg *paho.Publish) {
				ctx := context.Background()

				// extract trace id
				if mqttMsg.Properties != nil && mqttMsg.Properties.CorrelationData != nil {
					correlationMap := make(map[string]string)
					err := json.Unmarshal(mqttMsg.Properties.CorrelationData, &correlationMap)
					if err != nil {
						slog.Warn("failed to unmarshal correlation data", "error", err)
					} else {
						ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(correlationMap))
					}
				}

				// create span
				newCtx, span := l.tracer.Start(ctx,
					fmt.Sprintf("%s receive", getTopicPattern(mqttMsg.Topic)),
					trace.WithSpanKind(trace.SpanKindConsumer),
					trace.WithAttributes(
						semconv.MessagingSystem("mqtt"),
						semconv.MessagingConsumerID(clientId),
						semconv.MessagingMessagePayloadSizeBytes(len(mqttMsg.Payload)),
						semconv.MessagingOperationKey.String("receive"),
					))
				defer span.End()

				// determine charge station id
				topicParts := strings.Split(mqttMsg.Topic, "/")
				var chargeStationId = topicParts[len(topicParts)-1]

				// unmarshal the message
				var msg transport.Message
				err := json.Unmarshal(mqttMsg.Payload, &msg)
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, "unable to unmarshal message")
					slog.Warn("unable to unmarshal message", "err", err)
					return
				}

				// add additional span attributes
				version, _ := strings.CutPrefix(string(ocppVersion), "ocpp")
				span.SetAttributes(
					attribute.String("csId", chargeStationId),
					attribute.String("ocpp.version", version),
					attribute.String(fmt.Sprintf("%s.action", msg.MessageType), msg.Action),
					semconv.MessagingMessageConversationID(msg.MessageId),
				)

				if msg.MessageType == transport.MessageTypeCallError {
					span.SetAttributes(
						attribute.String(fmt.Sprintf("%s.code", msg.MessageType), string(msg.ErrorCode)),
						attribute.String(fmt.Sprintf("%s.description", msg.MessageType), msg.ErrorDescription))
				}

				// execute the handler
				handler.Handle(newCtx, chargeStationId, &msg)
			})
			readyCh <- struct{}{}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: clientId,
			Router:   mqttRouter,
		},
	})
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, errors.New("timeout waiting for mqtt connectionDetails setup")
	case <-readyCh:
		return conn, nil
	}
}

type connection struct {
	mqttConn *autopaho.ConnectionManager
}

func (c *connection) Disconnect(ctx context.Context) error {
	if c.mqttConn != nil {
		err := c.mqttConn.Disconnect(ctx)
		if err != nil {
			return err
		}
		c.mqttConn = nil
	}
	return nil
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

func getTopicPattern(topic string) string {
	parts := strings.Split(topic, "/")
	parts[len(parts)-1] = "#"
	return strings.Join(parts, "/")
}
