// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"go.opentelemetry.io/otel/trace"
	"net/url"
	"time"
)

type connection struct {
	mqttBrokerUrls        []*url.URL
	mqttPrefix            string
	mqttConnectTimeout    time.Duration
	mqttConnectRetryDelay time.Duration
	mqttKeepAliveInterval uint16
}

type Opt[T any] func(h *T)

func WithMqttBrokerUrl[T any](brokerUrl *url.URL) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Receiver:
			x.mqttBrokerUrls = append(x.mqttBrokerUrls, brokerUrl)
		case *Emitter:
			x.mqttBrokerUrls = append(x.mqttBrokerUrls, brokerUrl)
		}
	}
}

func WithMqttBrokerUrls[T any](brokerUrls []*url.URL) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Receiver:
			x.mqttBrokerUrls = brokerUrls
		case *Emitter:
			x.mqttBrokerUrls = brokerUrls
		}
	}
}

func WithMqttPrefix[T any](mqttPrefix string) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Receiver:
			x.mqttPrefix = mqttPrefix
		case *Emitter:
			x.mqttPrefix = mqttPrefix
		}
	}
}

func WithMqttConnectSettings[T any](mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval time.Duration) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Receiver:
			x.mqttConnectTimeout = mqttConnectTimeout
			x.mqttConnectRetryDelay = mqttConnectRetryDelay
			x.mqttKeepAliveInterval = uint16(mqttKeepAliveInterval.Round(time.Second).Seconds())
		case *Emitter:
			x.mqttConnectTimeout = mqttConnectTimeout
			x.mqttConnectRetryDelay = mqttConnectRetryDelay
			x.mqttKeepAliveInterval = uint16(mqttKeepAliveInterval.Round(time.Second).Seconds())
		}
	}
}

func WithOtelTracer[T any](tracer trace.Tracer) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Receiver:
			x.tracer = tracer
		case *Emitter:
			x.tracer = tracer
		}
	}
}

func WithMqttGroup[T Receiver](mqttGroup string) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Receiver:
			x.mqttGroup = mqttGroup
		}
	}
}

func WithRouter[T Receiver](router transport.Router) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Receiver:
			x.routers[router.GetOcppVersion()] = router
		}
	}
}

func WithEmitter[T Receiver](emitter transport.Emitter) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Receiver:
			x.emitter = emitter
		}
	}
}
