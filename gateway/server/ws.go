// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thoughtworks/maeve-csms/gateway/ocpp"
	"github.com/thoughtworks/maeve-csms/gateway/pipe"
	"github.com/thoughtworks/maeve-csms/gateway/registry"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
	"nhooyr.io/websocket"
)

type WebsocketHandler struct {
	mqttBrokerURLs        []*url.URL
	mqttTopicPrefix       string
	mqttConnectTimeout    time.Duration
	mqttConnectRetryDelay time.Duration
	mqttKeepAliveInterval uint16
	deviceRegistry        registry.DeviceRegistry
	orgNames              []string
	pipeOptions           []pipe.Opt
	trustProxyHeaders     bool
	tracer                trace.Tracer
}

type WebsocketOpt func(handler *WebsocketHandler)

func WithMqttBrokerUrl(brokerUrl *url.URL) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.mqttBrokerURLs = append(handler.mqttBrokerURLs, brokerUrl)
	}
}

func WithMqttBrokerUrlString(brokerUrl string) WebsocketOpt {
	u, err := url.Parse(brokerUrl)
	if err != nil {
		panic(err)
	}
	return func(handler *WebsocketHandler) {
		handler.mqttBrokerURLs = append(handler.mqttBrokerURLs, u)
	}
}

func WithMqttBrokerUrls(brokerUrls []*url.URL) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.mqttBrokerURLs = append(handler.mqttBrokerURLs, brokerUrls...)
	}
}

func WithMqttTopicPrefix(topicPrefix string) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.mqttTopicPrefix = topicPrefix
	}
}

func WithMqttConnectSettings(mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval time.Duration) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.mqttConnectTimeout = mqttConnectTimeout
		handler.mqttConnectRetryDelay = mqttConnectRetryDelay
		handler.mqttKeepAliveInterval = uint16(mqttKeepAliveInterval.Round(time.Second).Seconds())
	}
}

func WithDeviceRegistry(deviceRegistry registry.DeviceRegistry) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.deviceRegistry = deviceRegistry
	}
}

func WithOrgName(orgName string) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.orgNames = append(handler.orgNames, orgName)
	}
}

func WithOrgNames(orgNames []string) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.orgNames = append(handler.orgNames, orgNames...)
	}
}

func WithTrustProxyHeaders(trustProxyHeaders bool) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.trustProxyHeaders = trustProxyHeaders
	}
}

func WithPipeOption(pipeOption pipe.Opt) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.pipeOptions = append(handler.pipeOptions, pipeOption)
	}
}

func WithPipeOptions(pipeOption []pipe.Opt) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.pipeOptions = append(handler.pipeOptions, pipeOption...)
	}
}

func WithOtelTracer(tracer trace.Tracer) WebsocketOpt {
	return func(handler *WebsocketHandler) {
		handler.tracer = tracer
	}
}

func NewWebsocketHandler(opts ...WebsocketOpt) http.Handler {
	s := new(WebsocketHandler)

	for _, opt := range opts {
		opt(s)
	}

	ensureDefaults(s)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(TraceRequest(s.tracer))
	if s.trustProxyHeaders {
		r.Use(TLSOffload(s.deviceRegistry))
	}
	r.Handle("/ws/{id}", s)
	return r
}

func ensureDefaults(handler *WebsocketHandler) {
	if handler.mqttBrokerURLs == nil {
		u, err := url.Parse("mqtt://127.0.0.1:1883/")
		if err != nil {
			panic(err)
		}
		handler.mqttBrokerURLs = []*url.URL{u}
	}

	if handler.mqttConnectTimeout == 0 {
		handler.mqttConnectTimeout = 5 * time.Second
	}

	if handler.mqttConnectRetryDelay == 0 {
		handler.mqttConnectRetryDelay = 1 * time.Second
	}

	if handler.mqttKeepAliveInterval == 0 {
		handler.mqttKeepAliveInterval = 10
	}

	if handler.deviceRegistry == nil {
		panic("must provide device registry implementation")
	}

	if handler.tracer == nil {
		handler.tracer = trace.NewNoopTracerProvider().Tracer("")
	}
}

func (s *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("websocket connection received", "path", r.URL.Path, "method", r.Method)
	slog.Info("processing connection", "uri", r.RequestURI)

	span := trace.SpanFromContext(r.Context())

	clientId := chi.URLParam(r, "id")
	if clientId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("csId", clientId))

	cs, err := s.deviceRegistry.LookupChargeStation(clientId)
	if err != nil {
		span.SetStatus(codes.Error, "lookup charge station failed")
		span.RecordError(err)
		span.SetAttributes(semconv.HTTPStatusCode(http.StatusInternalServerError))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if cs == nil {
		span.SetStatus(codes.Error, "unknown charge station")
		span.RecordError(err)
		span.SetAttributes(semconv.HTTPStatusCode(http.StatusNotFound))
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	span.SetAttributes(attribute.Int("ocpp.security_profile", int(cs.SecurityProfile)))

	switch cs.SecurityProfile {
	case registry.UnsecuredTransportWithBasicAuth:
		if r.TLS != nil || !checkAuthorization(r.Context(), r, cs) {
			if r.TLS != nil {
				span.SetAttributes(attribute.String("auth.failure_reason", "tls for unsecured transport"))
			}
			span.SetStatus(codes.Error, "unauthorized")
			span.SetAttributes(semconv.HTTPStatusCode(http.StatusUnauthorized))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	case registry.TLSWithBasicAuth:
		if r.TLS == nil || !checkAuthorization(r.Context(), r, cs) {
			if r.TLS == nil {
				span.SetAttributes(attribute.String("auth.failure_reason", "no tls for secured transport"))
			}
			span.SetStatus(codes.Error, "unauthorized")
			span.SetAttributes(semconv.HTTPStatusCode(http.StatusUnauthorized))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	case registry.TLSWithClientSideCertificates:
		if r.TLS == nil || !checkCertificate(r.Context(), r, s.orgNames, cs) {
			if r.TLS == nil {
				span.SetAttributes(attribute.String("auth.failure_reason", "no tls for secured transport"))
			}
			span.SetStatus(codes.Error, "unauthorized")
			span.SetAttributes(semconv.HTTPStatusCode(http.StatusUnauthorized))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{Subprotocols: []string{"ocpp2.0.1", "ocpp1.6"}, InsecureSkipVerify: true})
	if err != nil {
		span.SetAttributes(attribute.String("websocket.accept_failure_reason", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	protocol := wsConn.Subprotocol()
	if protocol == "" {
		protocol = "ocpp2.0.1"
	}

	span.SetAttributes(attribute.String("ocpp.protocol", protocol))

	p := pipe.NewPipe(s.pipeOptions...)
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	mqttBrokerURLStrings := make([]string, len(s.mqttBrokerURLs))
	for i, u := range s.mqttBrokerURLs {
		mqttBrokerURLStrings[i] = u.String()
	}
	span.SetAttributes(attribute.StringSlice("mqtt.broker_urls", mqttBrokerURLStrings))

	var mqttConn *autopaho.ConnectionManager
	mqttConn, err = autopaho.NewConnection(ctx, autopaho.ClientConfig{
		BrokerUrls:        s.mqttBrokerURLs,
		KeepAlive:         s.mqttKeepAliveInterval,
		ConnectRetryDelay: s.mqttConnectRetryDelay,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			topicName := fmt.Sprintf("%s/out/%s/%s", s.mqttTopicPrefix, protocol, clientId)
			span.SetAttributes(attribute.String("mqtt.topic", topicName))
			_, err = manager.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					topicName: {},
				},
			})
			if err != nil {
				span.SetStatus(codes.Error, "subscribing to mqtt topic failed")
				span.RecordError(err)
				span.SetAttributes(semconv.HTTPStatusCode(http.StatusInternalServerError))

				_ = wsConn.Close(websocket.StatusProtocolError, http.StatusText(http.StatusInternalServerError))
				return
			}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: clientId,
			Router: paho.NewSingleHandlerRouter(func(mqttMsg *paho.Publish) {
				// route requests from the CSMS
				var msg pipe.GatewayMessage
				err := json.Unmarshal(mqttMsg.Payload, &msg)
				if err != nil {
					slog.Error("unmarshalling CSMS message", "err", err)
					return
				}

				correlationMap := make(map[string]string)
				err = json.Unmarshal(mqttMsg.Properties.CorrelationData, &correlationMap)
				if err != nil {
					slog.Warn("unmarshalling correlation map", "err", err)
				}
				requestContext := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(correlationMap))

				newCtx, span := s.tracer.Start(requestContext, fmt.Sprintf("%s/out/%s/# receive", s.mqttTopicPrefix, protocol),
					trace.WithSpanKind(trace.SpanKindConsumer),
					trace.WithAttributes(
						semconv.MessagingSystem("mqtt"),
						semconv.MessagingMessagePayloadSizeBytes(len(mqttMsg.Payload)),
						semconv.MessagingMessageConversationID(msg.MessageId),
						semconv.MessagingOperationKey.String("receive"),
						attribute.String("csId", clientId),
					))
				defer span.End()

				msg.Context = newCtx

				p.CSMSRx <- &msg
			}),
			OnServerDisconnect: func(disconnect *paho.Disconnect) {
				span.SetAttributes(attribute.String("mqtt.disconnect_reason", disconnect.Properties.ReasonString))
			},
		},
	})
	if err != nil {
		span.SetStatus(codes.Error, "connecting to mqtt")
		span.RecordError(err)
		span.SetAttributes(semconv.HTTPStatusCode(http.StatusInternalServerError))
		slog.Error("connecting to mqtt", "mqttBrokerURLs", s.mqttBrokerURLs, "err", err)
		_ = wsConn.Close(websocket.StatusProtocolError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer func() {
		err := mqttConn.Disconnect(context.Background())
		if err != nil {
			slog.Error("disconnecting from mqtt", "err", err)
		}
	}()
	err = mqttConn.AwaitConnection(ctx)
	if err != nil {
		slog.Error("waiting for mqtt", "mqttBrokerURLs", s.mqttBrokerURLs, "err", err)
		return
	}

	// we've finished connecting... complete this span so we get to see the details in the trace
	span.End()

	// listen on the CSMS Tx channel and publish those messages on the inbound topic
	goPublishToCSMS(ctx, s.tracer, p.CSMSTx, mqttConn, s.mqttTopicPrefix, protocol, clientId)

	// listen the CS Tx channel and write those messages to the websocket
	goWriteToChargeStation(ctx, s.tracer, p.ChargeStationTx, wsConn, protocol, clientId)

	// read from the websocket and send to the CS Rx channel (CS Tx used for error)
	readFromChargeStation(ctx, s.tracer, wsConn, p.ChargeStationRx, p.ChargeStationTx, protocol, clientId)
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "wss"
	}
	return "ws"
}

func checkAuthorization(ctx context.Context, r *http.Request, cs *registry.ChargeStation) bool {
	span := trace.SpanFromContext(ctx)

	username, password, ok := r.BasicAuth()
	if !ok {
		span.SetAttributes(attribute.String("auth.failure_reason", "no basic auth"))
		return false
	}
	if username != cs.ClientId {
		span.SetAttributes(attribute.String("auth.failure_reason", "invalid username"))
		return false
	}
	sha256pw := sha256.Sum256([]byte(password))
	b64sha256 := base64.StdEncoding.EncodeToString(sha256pw[:])
	result := b64sha256 == cs.Base64SHA256Password

	if !result {
		span.SetAttributes(attribute.String("auth.failure_reason", "invalid password"))
	}

	return result
}

func checkCertificate(ctx context.Context, r *http.Request, orgNames []string, cs *registry.ChargeStation) bool {
	span := trace.SpanFromContext(ctx)

	if len(r.TLS.PeerCertificates) == 0 {
		span.SetAttributes(attribute.String("auth.failure_reason", "no client certificate"))
		return false
	}

	leafCertificate := r.TLS.PeerCertificates[0]

	foundOrg := false

	span.SetAttributes(
		attribute.StringSlice("auth.organization", leafCertificate.Subject.Organization),
		attribute.String("auth.common_name", leafCertificate.Subject.CommonName))

	for _, org := range leafCertificate.Subject.Organization {

		if slices.Contains(orgNames, org) {
			foundOrg = true
			break
		}
	}
	if !foundOrg {
		span.SetAttributes(attribute.String("auth.failure_reason", "bad organization"))
		return false
	}

	//result := cs.ClientId == leafCertificate.Subject.CommonName
	//
	//if !result {
	//	span.SetAttributes(attribute.String("auth.failure_reason", "bad common name"))
	//}

	return foundOrg
}

func goPublishToCSMS(ctx context.Context, tracer trace.Tracer, csmsTx chan *pipe.GatewayMessage, mqttConn *autopaho.ConnectionManager, topicPrefix, protocol, clientId string) {
	go func() {
		for {
			select {
			case msg := <-csmsTx:
				data, err := json.Marshal(msg)
				if err != nil {
					slog.Error("marshaling message for publication", "err", err)
					continue
				}
				err = publish(msg.Context, tracer, mqttConn, topicPrefix, protocol, clientId, msg.MessageId, data)
				if err != nil {
					slog.Error("publishing message", "err", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func publish(ctx context.Context, tracer trace.Tracer, mqttConn *autopaho.ConnectionManager, topicPrefix, protocol, clientId, messageId string, data []byte) error {
	topic := fmt.Sprintf("%s/in/%s/%s", topicPrefix, protocol, clientId)

	newCtx, span := tracer.Start(ctx,
		fmt.Sprintf("%s/in/%s/# publish", topicPrefix, protocol),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.MessagingSystem("mqtt"),
			semconv.MessagingMessagePayloadSizeBytes(len(data)),
			semconv.MessagingOperationKey.String("publish"),
			semconv.MessagingMessageConversationID(messageId),
			attribute.String("csId", clientId),
		))
	defer span.End()

	correlationMap := make(map[string]string)
	otel.GetTextMapPropagator().Inject(newCtx, propagation.MapCarrier(correlationMap))

	correlationData, err := json.Marshal(correlationMap)
	if err != nil {
		slog.Warn("marshalling correlation map: %v", err)
	}

	_, err = mqttConn.Publish(newCtx, &paho.Publish{
		Topic:   topic,
		Payload: data,
		Properties: &paho.PublishProperties{
			ContentType:     "application/json",
			ResponseTopic:   fmt.Sprintf("%s/out/%s/%s", topicPrefix, protocol, clientId),
			CorrelationData: correlationData,
		},
	})

	return err
}

func goWriteToChargeStation(ctx context.Context, tracer trace.Tracer, chargeStationTx chan *pipe.GatewayMessage, wsConn *websocket.Conn, protocol, clientId string) {
	go func() {
		for {
			select {
			case msg := <-chargeStationTx:
				data, err := marshalGatewayMessageAsOcpp(msg)
				if err != nil {
					slog.Error("marshaling gateway message for charge station", "err", err)
					continue
				}
				err = write(msg.Context, tracer, wsConn, protocol, clientId, msg.MessageId, data)
				if err != nil {
					slog.Error("writing to charge station", "err", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func write(ctx context.Context, tracer trace.Tracer, wsConn *websocket.Conn, protocol, clientId, messageId string, data []byte) error {
	newCtx, span := tracer.Start(ctx, fmt.Sprintf("%s publish", protocol), trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			semconv.MessagingSystem("websocket"),
			semconv.MessagingMessagePayloadSizeBytes(len(data)),
			semconv.MessagingOperationKey.String("publish"),
			semconv.MessagingMessageConversationID(messageId),
			attribute.String("csId", clientId),
		))
	defer span.End()

	return wsConn.Write(newCtx, websocket.MessageText, data)
}

func readFromChargeStation(ctx context.Context, tracer trace.Tracer, wsConn *websocket.Conn, csRx, csTx chan *pipe.GatewayMessage, protocol, clientId string) {
	for {
		msg, err := read(ctx, tracer, wsConn, protocol, clientId)
		if err != nil {
			if msg != nil {
				slog.Warn("sending error message to client", "err", err)
				csTx <- msg
				continue
			}
			// connection closed
			break
		} else if msg != nil {
			csRx <- msg
		} else {
			// deadline exceeded
			break
		}
	}
}

var errClient = errors.New("client error")

func read(ctx context.Context, tracer trace.Tracer, wsConn *websocket.Conn, protocol, clientId string) (*pipe.GatewayMessage, error) {
	typ, b, err := wsConn.Read(context.Background())
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, nil
	} else if status := websocket.CloseStatus(err); status != -1 {
		slog.Info("connection closed with status", "status", status)
		return nil, err
	} else if errors.Is(err, io.EOF) {
		slog.Info("connection closed")
		return nil, err
	} else if err != nil {
		return nil, err
	}

	newCtx, span := tracer.Start(context.Background(), fmt.Sprintf("%s receive", protocol), trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			semconv.MessagingSystem("websocket"),
			semconv.MessagingOperationKey.String("receive"),
			semconv.MessagingMessagePayloadSizeBytes(len(b)),
			attribute.String("csId", clientId),
		),
		trace.WithLinks(trace.Link{
			SpanContext: trace.SpanContextFromContext(ctx),
		}))
	defer span.End()

	if typ != websocket.MessageText {
		msg := pipe.GatewayMessage{
			Context:          newCtx,
			MessageType:      ocpp.MessageTypeCallError,
			MessageId:        "-1",
			ErrorCode:        ocpp.ErrorRpcFrameworkError,
			ErrorDescription: "websocket message type is not text",
		}
		span.SetAttributes(semconv.MessagingMessageConversationID("-1"))
		span.SetStatus(codes.Error, "unmarshal ocpp message error")
		span.RecordError(err)

		return &msg, errClient
	}

	msg, err := unmarshalOcppAsGatewayMessage(b)
	if err != nil {
		msg := pipe.GatewayMessage{
			Context:          newCtx,
			MessageType:      ocpp.MessageTypeCallError,
			MessageId:        "-1",
			ErrorCode:        ocpp.ErrorRpcFrameworkError,
			ErrorDescription: err.Error(),
		}
		span.SetAttributes(semconv.MessagingMessageConversationID("-1"))
		span.SetStatus(codes.Error, "unmarshal ocpp message error")
		span.RecordError(err)
		return &msg, errClient
	}

	msg.Context = newCtx

	span.SetAttributes(semconv.MessagingMessageConversationID(msg.MessageId))

	return msg, nil
}

func marshalGatewayMessageAsOcpp(msg *pipe.GatewayMessage) ([]byte, error) {
	var err error
	ocppMsg := ocpp.Message{}
	ocppMsg.MessageTypeId = msg.MessageType
	ocppMsg.MessageId = msg.MessageId

	switch msg.MessageType {
	case ocpp.MessageTypeCall:
		ocppMsg.Data = make([]json.RawMessage, 2)
		ocppMsg.Data[0], err = json.Marshal(msg.Action)
		if err != nil {
			return nil, err
		}
		ocppMsg.Data[1] = msg.RequestPayload
	case ocpp.MessageTypeCallResult:
		ocppMsg.Data = make([]json.RawMessage, 1)
		ocppMsg.Data[0] = msg.ResponsePayload
	case ocpp.MessageTypeCallError:
		ocppMsg.Data = make([]json.RawMessage, 3)
		ocppMsg.Data[0], err = json.Marshal(msg.ErrorCode)
		if err != nil {
			return nil, err
		}
		ocppMsg.Data[1], err = json.Marshal(msg.ErrorDescription)
		if err != nil {
			return nil, err
		}
		ocppMsg.Data[2] = json.RawMessage("{}")
	}

	return json.Marshal(ocppMsg)
}

func unmarshalOcppAsGatewayMessage(b []byte) (*pipe.GatewayMessage, error) {
	ocppMsg := ocpp.Message{}
	err := json.Unmarshal(b, &ocppMsg)
	if err != nil {
		return nil, err
	}
	msg := pipe.GatewayMessage{
		MessageType: ocppMsg.MessageTypeId,
		MessageId:   ocppMsg.MessageId,
	}
	switch ocppMsg.MessageTypeId {
	case ocpp.MessageTypeCall:
		err = json.Unmarshal(ocppMsg.Data[0], &msg.Action)
		if err != nil {
			return nil, err
		}
		msg.RequestPayload = ocppMsg.Data[1]
	case ocpp.MessageTypeCallResult:
		msg.ResponsePayload = ocppMsg.Data[0]
	case ocpp.MessageTypeCallError:
		err = json.Unmarshal(ocppMsg.Data[0], &msg.ErrorCode)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(ocppMsg.Data[1], &msg.ErrorDescription)
		if err != nil {
			return nil, err
		}
	}
	return &msg, nil
}
