// SPDX-License-Identifier: Apache-2.0

package server_test

import (
	"context"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/thoughtworks/maeve-csms/gateway/server"
	"net/url"
	"testing"
	"time"
)

func TestNewBroker(t *testing.T) {
	broker, addr := server.NewBroker(t)
	defer func() {
		err := broker.Close()
		if err != nil {
			t.Errorf("closing broker: %v", err)
		}
	}()

	err := broker.Serve()
	if err != nil {
		t.Fatalf("starting broker: %v", err)
	}

	doneCh := make(chan struct{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{addr},
		KeepAlive:         10,
		ConnectRetryDelay: 10,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err = manager.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					"test/#": {},
				},
			})
			if err != nil {
				t.Fatalf("subscribing to messages: %v", err)
			}

			err = broker.Publish("test/123", []byte("test data"), false, 0)
			if err != nil {
				t.Errorf("publising message: %v", err)
			}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "cs1",
			Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
				t.Logf("received message: %v", m)
				doneCh <- struct{}{}
			}),
		},
	})
	if err != nil {
		t.Fatalf("connecting to broker: %v", err)
	}

	select {
	case <-doneCh:
		return
	case <-ctx.Done():
		t.Fatal("failed to receive message")
	}
}
