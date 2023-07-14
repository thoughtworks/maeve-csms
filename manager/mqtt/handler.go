// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
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

	mqttRouter := paho.NewStandardRouter()
	mqttRouter.RegisterHandler(mqttV16Topic, NewGatewayMessageHandler(context.Background(), v16Router, v16Emitter, h.schemaFS))
	mqttRouter.RegisterHandler(mqttV201Topic, NewGatewayMessageHandler(context.Background(), v201Router, v201Emitter, h.schemaFS))

	var mqttConn *autopaho.ConnectionManager
	mqttConn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        h.mqttBrokerUrls,
		KeepAlive:         h.mqttKeepAliveInterval,
		ConnectRetryDelay: h.mqttConnectRetryDelay,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			v16Emitter.emitter = NewMqttEmitter(mqttConn, h.mqttPrefix, "ocpp1.6")
			v201Emitter.emitter = NewMqttEmitter(mqttConn, h.mqttPrefix, "ocpp2.0.1")

			_, err := mqttConn.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					mqttV16Topic:  {},
					mqttV201Topic: {},
				},
			})
			if err != nil {
				errCh <- err
			}

			log.Printf("subscribed to gateway+ocpp1.6 on %+v/%s", h.mqttBrokerUrls, mqttV16Topic)
			log.Printf("subscribed to gateway+ocpp2.0.1 on %+v/%s", h.mqttBrokerUrls, mqttV201Topic)

			readyCh <- struct{}{}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: fmt.Sprintf("%s-%s", h.mqttGroup, randSeq(5)),
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

func NewGatewayMessageHandler(ctx context.Context, router *Router, emitter Emitter, schemaFS fs.FS) func(mqttMsg *paho.Publish) {
	return func(mqttMsg *paho.Publish) {
		topicParts := strings.Split(mqttMsg.Topic, "/")
		chargeStationId := topicParts[len(topicParts)-1]
		var msg Message
		err := json.Unmarshal(mqttMsg.Payload, &msg)
		if err != nil {
			errMsg := NewErrorMessage("", "-1", ErrorInternalError, err)
			err = emitter.Emit(ctx, chargeStationId, errMsg)
			if err != nil {
				log.Printf("unable to emit error message: %v", err)
			}
		}
		err = router.Route(ctx, chargeStationId, msg, emitter, schemaFS)
		if err != nil {
			log.Printf("ERROR: %s - %s: %v", chargeStationId, msg.Action, err)
			var mqttError *Error
			var errMsg *Message
			if errors.As(err, &mqttError) {
				errMsg = NewErrorMessage(msg.Action, msg.MessageId, mqttError.ErrorCode, mqttError.wrappedError)
			} else {
				errMsg = NewErrorMessage(msg.Action, msg.MessageId, ErrorInternalError, err)
			}
			err = emitter.Emit(ctx, chargeStationId, errMsg)
			if err != nil {
				log.Printf("unable to emit error message: %v", err)
			}
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
