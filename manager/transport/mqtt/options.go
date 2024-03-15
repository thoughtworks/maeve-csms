// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"go.opentelemetry.io/otel/trace"
	"net/url"
	"time"
)

type connectionDetails struct {
	mqttBrokerUrls        []*url.URL
	mqttPrefix            string
	mqttConnectTimeout    time.Duration
	mqttConnectRetryDelay time.Duration
	mqttKeepAliveInterval uint16
}

type Opt[T any] func(h *T)

func WithMqttBrokerUrl[T Emitter | Listener](brokerUrl *url.URL) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttBrokerUrls = append(x.mqttBrokerUrls, brokerUrl)
		case *Listener:
			x.mqttBrokerUrls = append(x.mqttBrokerUrls, brokerUrl)
		}
	}
}

func WithMqttBrokerUrls[T Emitter | Listener](brokerUrls []*url.URL) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttBrokerUrls = brokerUrls
		case *Listener:
			x.mqttBrokerUrls = brokerUrls
		}
	}
}

func WithMqttPrefix[T Emitter | Listener](mqttPrefix string) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttPrefix = mqttPrefix
		case *Listener:
			x.mqttPrefix = mqttPrefix
		}
	}
}

func WithMqttConnectSettings[T Emitter | Listener](mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval time.Duration) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.mqttConnectTimeout = mqttConnectTimeout
			x.mqttConnectRetryDelay = mqttConnectRetryDelay
			x.mqttKeepAliveInterval = uint16(mqttKeepAliveInterval.Round(time.Second).Seconds())
		case *Listener:
			x.mqttConnectTimeout = mqttConnectTimeout
			x.mqttConnectRetryDelay = mqttConnectRetryDelay
			x.mqttKeepAliveInterval = uint16(mqttKeepAliveInterval.Round(time.Second).Seconds())
		}
	}
}

func WithOtelTracer[T Emitter | Listener](tracer trace.Tracer) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Emitter:
			x.tracer = tracer
		case *Listener:
			x.tracer = tracer
		}
	}
}

func WithMqttGroup[T Listener](mqttGroup string) Opt[T] {
	return func(h *T) {
		switch x := any(h).(type) {
		case *Listener:
			x.mqttGroup = mqttGroup
		}
	}
}
