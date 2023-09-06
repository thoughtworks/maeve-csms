package mqtt

import (
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
	"golang.org/x/net/context"
	"net/url"
	"reflect"
	"time"
)

type Sender struct {
	mqttBrokerUrls        []*url.URL
	mqttPrefix            string
	mqttGroup             string
	mqttConnectTimeout    time.Duration
	mqttConnectRetryDelay time.Duration
	mqttKeepAliveInterval uint16
	tracer                trace.Tracer

	V16CallMaker BasicCallMaker
}

func NewSender(mqttBrokerUrls []*url.URL,
	mqttPrefix,
	mqttGroup string,
	mqttConnectTimeout, mqttConnectRetryDelay time.Duration,
	mqttKeepAliveInterval uint16,
	tracer trace.Tracer,
) *Sender {
	return &Sender{
		mqttBrokerUrls:        mqttBrokerUrls,
		mqttPrefix:            mqttPrefix,
		mqttGroup:             mqttGroup,
		mqttConnectTimeout:    mqttConnectTimeout,
		mqttConnectRetryDelay: mqttConnectRetryDelay,
		mqttKeepAliveInterval: mqttKeepAliveInterval,
		tracer:                tracer,
	}
}

func (s *Sender) Connect(errChan chan error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	readyCh := make(chan struct{})

	var v16Emitter Emitter
	v16Emitter = &ProxyEmitter{}

	var mqttConn *autopaho.ConnectionManager
	mqttConn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        s.mqttBrokerUrls,
		KeepAlive:         s.mqttKeepAliveInterval,
		ConnectRetryDelay: s.mqttConnectRetryDelay,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			v16Emitter = NewMqttEmitter(mqttConn, s.mqttPrefix, "ocpp1.6", s.tracer)
			readyCh <- struct{}{}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: fmt.Sprintf("%s-%s", "manager", randSeq(5)),
			Router:   paho.NewStandardRouter(),
		},
	})
	if err != nil {
		slog.Error("error setting up mqttConn", "err", err)
		errChan <- err
	}

	select {
	case <-ctx.Done():
		slog.Error("timed out waiting for mqtt connection setup")
		errChan <- err
	case <-readyCh:
		// do nothing
		slog.Info("mqtt connection ready", "emitter", v16Emitter)
	}

	s.V16CallMaker = BasicCallMaker{
		E: v16Emitter,
		Actions: map[reflect.Type]string{
			reflect.TypeOf(&ocpp16.RemoteStartTransactionJson{}): "RemoteStartTransaction",
		},
	}
}
