// SPDX-License-Identifier: Apache-2.0

package server_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/gateway/ocpp"
	"github.com/thoughtworks/maeve-csms/gateway/pipe"
	"github.com/thoughtworks/maeve-csms/gateway/registry"
	"github.com/thoughtworks/maeve-csms/gateway/server"
	"math"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"nhooyr.io/websocket"
	"testing"
	"time"
)

func TestWebSocketHandler(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	broker, addr := server.NewBroker(t)
	err := broker.Serve()
	if err != nil {
		t.Fatalf("starting broker: %v", err)
	}
	defer func() {
		err := broker.Close()
		if err != nil {
			t.Logf("WARN: broker close: %v", err)
		}
	}()

	// simulate manager connection
	client, err := autopaho.NewConnection(ctx, autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{addr},
		KeepAlive:         10,
		ConnectRetryDelay: 2 * time.Second,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err = manager.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					"cs/in/ocpp2.0.1/+": {},
				},
			})
			require.NoError(t, err)
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "test",
			Router: paho.NewSingleHandlerRouter(func(publish *paho.Publish) {
				var reqMsg pipe.GatewayMessage
				err := json.Unmarshal(publish.Payload, &reqMsg)
				require.NoError(t, err)

				respMsg := pipe.GatewayMessage{
					MessageType:     ocpp.MessageTypeCallResult,
					MessageId:       reqMsg.MessageId,
					ResponsePayload: reqMsg.RequestPayload,
				}

				b, err := json.Marshal(respMsg)
				require.NoError(t, err)
				err = broker.Publish(publish.Properties.ResponseTopic, b, false, 0)
				require.NoError(t, err)
			}),
		},
	})
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			t.Logf("WARN: mqtt client disconnect: %v", err)
		}
	}()

	err = client.AwaitConnection(ctx)
	require.NoError(t, err)

	cs := &registry.ChargeStation{
		ClientId:             "cs1",
		SecurityProfile:      registry.UnsecuredTransportWithBasicAuth,
		Base64SHA256Password: "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg=", // password
	}

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations["cs1"] = cs

	srv := httptest.NewServer(server.NewWebsocketHandler(
		server.WithMqttBrokerUrl(addr),
		server.WithMqttTopicPrefix("cs"),
		server.WithDeviceRegistry(mockRegistry),
		server.WithMqttConnectSettings(15*time.Second, 15*time.Second, 5*time.Second)))
	defer srv.Close()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cs.ClientId, "password")))
	dialOptions := &websocket.DialOptions{
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
		HTTPHeader: http.Header{
			"authorization": []string{authHeader},
		},
	}

	conn, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/cs1", srv.URL), dialOptions)
	if err != nil {
		t.Fatalf("dialing server: %v", err)
	}
	defer func() {
		err := conn.Close(websocket.StatusNormalClosure, "OK")
		if err != nil {
			t.Logf("WARN: closing websocket connection: %v", err)
		}
	}()

	if conn.Subprotocol() != "ocpp2.0.1" {
		t.Errorf("subprotocol: want %s, got %s", "ocpp2.0.1", conn.Subprotocol())
	}

	call := ocpp.Message{
		MessageTypeId: ocpp.MessageTypeCall,
		MessageId:     "1",
		Data: []json.RawMessage{
			json.RawMessage(`"EchoRequest"`),
			json.RawMessage(`"Payload"`),
		},
	}
	data, err := json.Marshal(call)
	if err != nil {
		t.Fatalf("encoding call: %v", err)
	}

	err = conn.Write(ctx, websocket.MessageText, data)
	if err != nil {
		t.Fatalf("sending message: %v", err)
	}

	typ, b, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("receiving message: %v", err)
	}
	if typ != websocket.MessageText {
		t.Errorf("message type: want %d, got %d", websocket.MessageText, typ)
	}
	var msg ocpp.Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		t.Errorf("unmarshaling message: %v", err)
	}

	doneCh := make(chan struct{}, 1)
	if msg.MessageTypeId == ocpp.MessageTypeCallResult {
		if string(msg.Data[0]) != `"Payload"` {
			t.Fatalf("call result: want %s, got %s", "Payload", string(msg.Data[0]))
		}
		doneCh <- struct{}{}
	} else if msg.MessageTypeId == ocpp.MessageTypeCallError {
		t.Fatalf("received error: %s", msg.Data[1])
	}

	select {
	case <-doneCh:
		// do nothing
	case <-ctx.Done():
		t.Fatal("timeout waiting for test to complete")
	}
}

func TestConnectionFromUnknownChargeStation(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockRegistry := registry.NewMockRegistry()

	srv := httptest.NewServer(server.NewWebsocketHandler(server.WithDeviceRegistry(mockRegistry)))
	defer srv.Close()

	dialOptions := &websocket.DialOptions{
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
	}

	conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/unknownCS1", srv.URL), dialOptions)
	if err == nil {
		_ = conn.Close(websocket.StatusGoingAway, "Shutdown")
		t.Fatalf("expected error dialing CSMS")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status code: want %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestHttpConnectionWithBasicAuth(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := &registry.ChargeStation{
		ClientId:             "basicAuthCS1",
		SecurityProfile:      registry.UnsecuredTransportWithBasicAuth,
		Base64SHA256Password: "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg=", // password,
	}

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations[cs.ClientId] = cs

	srv := httptest.NewServer(server.NewWebsocketHandler(server.WithDeviceRegistry(mockRegistry)))
	defer srv.Close()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cs.ClientId, "password")))
	dialOptions := &websocket.DialOptions{
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
		HTTPHeader: http.Header{
			"authorization": []string{authHeader},
		},
	}

	conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
	if err != nil {
		t.Fatalf("dialing CSMS: %v", err)
	}
	defer func() {
		err := conn.Close(websocket.StatusGoingAway, "Shutdown")
		if err != nil {
			t.Logf("WARN: websocket close: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("status code: want %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}
}

func TestHttpConnectionWithBasicAuthWrongPassword(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := &registry.ChargeStation{
		ClientId:             "basicAuthCS1WrongPassword",
		SecurityProfile:      registry.UnsecuredTransportWithBasicAuth,
		Base64SHA256Password: "bPYV1byqx3g1Ko8fM2DSPwLzTsGC4lmJf9bOSF14cNQ=", // password2,
	}

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations[cs.ClientId] = cs

	srv := httptest.NewServer(server.NewWebsocketHandler(server.WithDeviceRegistry(mockRegistry)))
	defer srv.Close()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cs.ClientId, "password")))
	dialOptions := &websocket.DialOptions{
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
		HTTPHeader: http.Header{
			"authorization": []string{authHeader},
		},
	}

	conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
	if err == nil {
		_ = conn.Close(websocket.StatusGoingAway, "Shutdown")
		t.Fatalf("expected error dialing CSMS")
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status code: want %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestHttpConnectionWithBasicAuthWrongUsername(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := &registry.ChargeStation{
		ClientId:             "basicAuthCS1WrongUsername",
		SecurityProfile:      registry.UnsecuredTransportWithBasicAuth,
		Base64SHA256Password: "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg=", // password,
	}

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations[cs.ClientId] = cs

	srv := httptest.NewServer(server.NewWebsocketHandler(server.WithDeviceRegistry(mockRegistry)))
	defer srv.Close()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "wrong", "password")))
	dialOptions := &websocket.DialOptions{
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
		HTTPHeader: http.Header{
			"authorization": []string{authHeader},
		},
	}

	conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
	if err == nil {
		_ = conn.Close(websocket.StatusGoingAway, "Shutdown")
		t.Fatalf("expected error dialing CSMS")
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status code: want %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestHttpConnectionWithBasicAuthInvalidUsernameAllowed(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := &registry.ChargeStation{
		ClientId:               "basicAuthCS1WrongUsername",
		SecurityProfile:        registry.UnsecuredTransportWithBasicAuth,
		Base64SHA256Password:   "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg=", // password,
		InvalidUsernameAllowed: true,
	}

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations[cs.ClientId] = cs

	srv := httptest.NewServer(server.NewWebsocketHandler(server.WithDeviceRegistry(mockRegistry)))
	defer srv.Close()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "wrong", "password")))
	dialOptions := &websocket.DialOptions{
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
		HTTPHeader: http.Header{
			"authorization": []string{authHeader},
		},
	}

	conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
	if err != nil {
		t.Fatalf("dialing CSMS: %v", err)
	}
	defer func() {
		err := conn.Close(websocket.StatusGoingAway, "Shutdown")
		if err != nil {
			t.Logf("WARN: websocket close: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("status code: want %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}
}

func TestTlsConnectionWithBasicAuth(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := &registry.ChargeStation{
		ClientId:             "tlsWithBasicAuthCS1",
		SecurityProfile:      registry.TLSWithBasicAuth,
		Base64SHA256Password: "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg=", // password,
	}

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations[cs.ClientId] = cs

	srv := httptest.NewTLSServer(server.NewWebsocketHandler(server.WithDeviceRegistry(mockRegistry)))
	defer srv.Close()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cs.ClientId, "password")))
	dialOptions := &websocket.DialOptions{
		HTTPClient:   srv.Client(),
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
		HTTPHeader: http.Header{
			"authorization": []string{authHeader},
		},
	}

	conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
	if err != nil {
		t.Fatalf("dialing CSMS: %v", err)
	}
	defer func() {
		err := conn.Close(websocket.StatusGoingAway, "Shutdown")
		if err != nil {
			t.Logf("WARN: websocket close: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("status code: want %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}
}

func TestTlsConnectionWithBasicAuthWrongPassword(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := &registry.ChargeStation{
		ClientId:             "tlsWithBasicAuthCS1WrongPassword",
		SecurityProfile:      registry.TLSWithBasicAuth,
		Base64SHA256Password: "bPYV1byqx3g1Ko8fM2DSPwLzTsGC4lmJf9bOSF14cNQ=", // password2,
	}

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations[cs.ClientId] = cs

	srv := httptest.NewServer(server.NewWebsocketHandler(server.WithDeviceRegistry(mockRegistry)))
	defer srv.Close()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cs.ClientId, "password")))
	dialOptions := &websocket.DialOptions{
		HTTPClient:   srv.Client(),
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
		HTTPHeader: http.Header{
			"authorization": []string{authHeader},
		},
	}

	conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
	if err == nil {
		_ = conn.Close(websocket.StatusGoingAway, "Shutdown")
		t.Fatalf("expected error dialing CSMS")
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status code: want %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestTlsConnectionWithCertificateAuth(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := &registry.ChargeStation{
		ClientId:        "tlsWithCertificateAuth",
		SecurityProfile: registry.TLSWithClientSideCertificates,
	}

	caCert, _, clientCert, clientKeyPair := createTestKeyPairs(t, cs.ClientId)

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations[cs.ClientId] = cs

	srv := httptest.NewUnstartedServer(server.NewWebsocketHandler(
		server.WithDeviceRegistry(mockRegistry),
		server.WithOrgName("Thoughtworks")))
	clientCAPool := x509.NewCertPool()
	clientCAPool.AddCert(caCert)
	srv.TLS = &tls.Config{
		ClientCAs:  clientCAPool,
		ClientAuth: tls.VerifyClientCertIfGiven,
	}
	srv.StartTLS()
	defer srv.Close()

	rootCAPool := x509.NewCertPool()
	rootCAPool.AddCert(srv.Certificate())

	clientTlsCert := tls.Certificate{
		Certificate: [][]byte{clientCert.Raw, caCert.Raw},
		PrivateKey:  clientKeyPair,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      rootCAPool,
				Certificates: []tls.Certificate{clientTlsCert},
			},
		},
	}
	dialOptions := &websocket.DialOptions{
		HTTPClient:   httpClient,
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
	}

	conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
	if err != nil {
		t.Fatalf("dialing CSMS: %v", err)
	}
	defer func() {
		err := conn.Close(websocket.StatusGoingAway, "Shutdown")
		if err != nil {
			t.Logf("WARN: websocket close: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("status code: want %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}
}

func TestTlsConnectionWithCertificateAuthUntrustedRoot(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := &registry.ChargeStation{
		ClientId:        "tlsWithCertificateAuthUntrustedRoot",
		SecurityProfile: registry.TLSWithClientSideCertificates,
	}

	caCert, _, clientCert, clientKeyPair := createTestKeyPairs(t, cs.ClientId)

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations[cs.ClientId] = cs

	srv := httptest.NewUnstartedServer(server.NewWebsocketHandler(
		server.WithDeviceRegistry(mockRegistry),
		server.WithOrgName("Thoughtworks")))
	srv.TLS = &tls.Config{
		ClientAuth: tls.VerifyClientCertIfGiven,
	}
	srv.StartTLS()
	defer srv.Close()

	rootCAPool := x509.NewCertPool()
	rootCAPool.AddCert(srv.Certificate())

	clientTlsCert := tls.Certificate{
		Certificate: [][]byte{clientCert.Raw, caCert.Raw},
		PrivateKey:  clientKeyPair,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      rootCAPool,
				Certificates: []tls.Certificate{clientTlsCert},
			},
		},
	}
	dialOptions := &websocket.DialOptions{
		HTTPClient:   httpClient,
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
	}

	conn, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
	if err == nil {
		_ = conn.Close(websocket.StatusGoingAway, "Shutdown")
		t.Fatalf("expected error dialing CSMS")
	}
}

func TestTlsConnectionWithCertificateAuthExpired(t *testing.T) {
	//defer goleak.VerifyNone(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cs := &registry.ChargeStation{
		ClientId:        "tlsWithCertificateAuthExpired",
		SecurityProfile: registry.TLSWithClientSideCertificates,
	}

	caCert, _, clientCert, clientKeyPair := createTestKeyPairs(t, cs.ClientId)

	mockRegistry := registry.NewMockRegistry()
	mockRegistry.ChargeStations[cs.ClientId] = cs

	srv := httptest.NewUnstartedServer(server.NewWebsocketHandler(
		server.WithDeviceRegistry(mockRegistry),
		server.WithOrgName("Thoughtworks")))
	clientCAPool := x509.NewCertPool()
	clientCAPool.AddCert(caCert)
	srv.TLS = &tls.Config{
		Time:       func() time.Time { return time.Now().Add(10 * time.Minute) },
		ClientCAs:  clientCAPool,
		ClientAuth: tls.VerifyClientCertIfGiven,
	}
	srv.StartTLS()
	defer srv.Close()

	rootCAPool := x509.NewCertPool()
	rootCAPool.AddCert(srv.Certificate())

	clientTlsCert := tls.Certificate{
		Certificate: [][]byte{clientCert.Raw, caCert.Raw},
		PrivateKey:  clientKeyPair,
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      rootCAPool,
				Certificates: []tls.Certificate{clientTlsCert},
			},
		},
	}
	dialOptions := &websocket.DialOptions{
		HTTPClient:   httpClient,
		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
	}

	_, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
	if err == nil {
		t.Fatalf("expected error dialing CSMS")
	}
}

//func TestTlsConnectionWithCertificateAuthMismatchedName(t *testing.T) {
//	//defer goleak.VerifyNone(t)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	cs := &registry.ChargeStation{
//		ClientId:        "tlsWithCertificateAuthMismatchedName",
//		SecurityProfile: registry.TLSWithClientSideCertificates,
//	}
//
//	caCert, _, clientCert, clientKeyPair := createTestKeyPairs(t, cs.ClientId+"Wrong")
//
//	mockRegistry := registry.NewMockRegistry()
//	mockRegistry.ChargeStations[cs.ClientId] = cs
//
//	srv := httptest.NewUnstartedServer(server.NewWebsocketHandler(
//		server.WithDeviceRegistry(mockRegistry),
//		server.WithOrgName("Thoughtworks")))
//	clientCAPool := x509.NewCertPool()
//	clientCAPool.AddCert(caCert)
//	srv.TLS = &tls.Config{
//		ClientCAs:  clientCAPool,
//		ClientAuth: tls.VerifyClientCertIfGiven,
//	}
//	srv.StartTLS()
//	defer srv.Close()
//
//	rootCAPool := x509.NewCertPool()
//	rootCAPool.AddCert(srv.Certificate())
//
//	clientTlsCert := tls.Certificate{
//		Certificate: [][]byte{clientCert.Raw, caCert.Raw},
//		PrivateKey:  clientKeyPair,
//	}
//
//	httpClient := &http.Client{
//		Transport: &http.Transport{
//			TLSClientConfig: &tls.Config{
//				RootCAs:      rootCAPool,
//				Certificates: []tls.Certificate{clientTlsCert},
//			},
//		},
//	}
//	dialOptions := &websocket.DialOptions{
//		HTTPClient:   httpClient,
//		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
//	}
//
//	_, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
//	if err == nil {
//		t.Fatalf("expected error dialing CSMS")
//	}
//	if resp.StatusCode != http.StatusUnauthorized {
//		t.Fatalf("status code: want %d, got %d", http.StatusUnauthorized, resp.StatusCode)
//	}
//}

//func TestConnectionWhenUnableToConnectToMqtt(t *testing.T) {
//	//defer goleak.VerifyNone(t)
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	cs := &registry.ChargeStation{
//		ClientId:             "basicAuthCS1",
//		SecurityProfile:      registry.UnsecuredTransportWithBasicAuth,
//		Base64SHA256Password: "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg=", // password,
//	}
//
//	mockRegistry := registry.NewMockRegistry()
//	mockRegistry.ChargeStations[cs.ClientId] = cs
//
//	srv := httptest.NewServer(server.NewWebsocketHandler(
//		server.WithDeviceRegistry(mockRegistry),
//		server.WithMqttBrokerUrlString("mqtt://fizzbuz.example.com:1833/"),
//		server.WithMqttConnectSettings(200*time.Millisecond, 100*time.Millisecond, 0)))
//	defer srv.Close()
//
//	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cs.ClientId, "password")))
//	dialOptions := &websocket.DialOptions{
//		Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
//		HTTPHeader: http.Header{
//			"authorization": []string{authHeader},
//		},
//	}
//
//	conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
//	if err != nil {
//		t.Fatalf("dialing CSMS: %v", err)
//	}
//	if resp.StatusCode != http.StatusSwitchingProtocols {
//		t.Fatalf("status code: want %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
//	}
//
//	call := ocpp.Call{
//		MessageId: "1",
//		Action:    "EchoRequest",
//		Payload:   json.RawMessage("\"Payload\""),
//	}
//	data, err := json.Marshal(call)
//	if err != nil {
//		t.Fatalf("encoding call: %v", err)
//	}
//
//	err = conn.Write(ctx, websocket.MessageText, data)
//	if err != nil {
//		t.Fatalf("writing to websocket: %v", err)
//	}
//
//	_, _, err = conn.Read(ctx)
//	if status := websocket.CloseStatus(err); status != websocket.StatusProtocolError {
//		if status == -1 {
//			_ = conn.Close(websocket.StatusGoingAway, "Shutdown")
//		}
//		t.Errorf("reading from websocket: want status %d got status %d", websocket.StatusProtocolError, status)
//	}
//}

//	func TestConnectionWhenChargeStationCallReceivedAndNoCSMSResponse(t *testing.T) {
//		defer goleak.VerifyNone(t)
//
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//
//		broker, addr := server.NewBroker(t)
//		err := broker.Serve()
//		if err != nil {
//			t.Fatalf("starting broker: %v", err)
//		}
//		defer func() {
//			err := broker.Close()
//			if err != nil {
//				t.Logf("WARN: broker close: %v", err)
//			}
//		}()
//
//		cs := &registry.ChargeStation{
//			ClientId:             "basicAuthCS1",
//			SecurityProfile:      registry.UnsecuredTransportWithBasicAuth,
//			Base64SHA256Password: "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg=", // password,
//		}
//
//		mockRegistry := registry.NewMockRegistry()
//		mockRegistry.ChargeStations[cs.ClientId] = cs
//
//		srv := httptest.NewServer(server.NewWebsocketHandler(
//			server.WithDeviceRegistry(mockRegistry),
//			server.WithMqttBrokerUrl(addr),
//			server.WithPipeOption(pipe.WithResponseTimeout(200*time.Millisecond))))
//		defer srv.Close()
//
//		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cs.ClientId, "password")))
//		dialOptions := &websocket.DialOptions{
//			Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
//			HTTPHeader: http.Header{
//				"authorization": []string{authHeader},
//			},
//		}
//
//		conn, resp, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/%s", srv.URL, cs.ClientId), dialOptions)
//		if err != nil {
//			t.Fatalf("dialing CSMS: %v", err)
//		}
//		defer func() {
//			err := conn.Close(websocket.StatusGoingAway, "Shutdown")
//			if err != nil {
//				t.Logf("websocket close: %v", err)
//			}
//		}()
//		if resp.StatusCode != http.StatusSwitchingProtocols {
//			t.Fatalf("status code: want %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
//		}
//
//		call := ocpp.Call{
//			MessageId: "1",
//			Action:    "EchoRequest",
//			Payload:   json.RawMessage("\"Payload\""),
//		}
//		data, err := json.Marshal(call)
//		if err != nil {
//			t.Fatalf("encoding call: %v", err)
//		}
//
//		err = conn.Write(ctx, websocket.MessageText, data)
//		if err != nil {
//			t.Fatalf("writing to websocket: %v", err)
//		}
//
//		typ, b, err := conn.Read(ctx)
//		if status := websocket.CloseStatus(err); status != -1 {
//			t.Errorf("reading from websocket: want status %d got status %d", -1, status)
//		}
//
//		if typ != websocket.MessageText {
//			t.Errorf("message type: want %d got %d", websocket.MessageText, typ)
//		}
//
//		var msg ocpp.Message
//		err = json.Unmarshal(b, &msg)
//		if err != nil {
//			t.Fatalf("unmarshalling response: %v", err)
//		}
//
//		if msg.MessageTypeId != ocpp.MessageTypeCallError {
//			t.Errorf("message type id: want %d got %d", ocpp.MessageTypeCallError, msg.MessageTypeId)
//		}
//	}
//
//	func TestConnectionShutsDownWhenChargeStationDisconnectsMidRequest(t *testing.T) {
//		defer goleak.VerifyNone(t)
//
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//
//		broker, addr := server.NewBroker(t)
//		err := broker.Serve()
//		if err != nil {
//			t.Fatalf("starting broker: %v", err)
//		}
//		defer func() {
//			err := broker.Close()
//			if err != nil {
//				t.Logf("WARN: broker close: %v", err)
//			}
//		}()
//
//		client := mqtttest.NewClient(t, ctx, addr, "test", func(mqttMsg *paho.Publish) {
//			var msg ocpp.Message
//			err := json.Unmarshal(mqttMsg.Payload, &msg)
//			if err != nil {
//				t.Logf("ERROR: %v", err)
//			}
//
//			time.Sleep(200 * time.Millisecond)
//
//			res := ocpp.CallResult{
//				MessageId: msg.MessageId,
//				Payload:   msg.Data[1],
//			}
//
//			data, err := json.Marshal(res)
//			if err != nil {
//				t.Logf("marshaling error: %v", err)
//			}
//
//			err = broker.Publish(mqttMsg.Properties.ResponseTopic, data, false, 0)
//			if err != nil {
//				t.Logf("sending response: %v", err)
//			}
//		})
//		defer func() {
//			err := client.Disconnect(ctx)
//			if err != nil {
//				t.Logf("WARN: mqtt client disconnect: %v", err)
//			}
//		}()
//
//		_, err = client.Subscribe(ctx, &paho.Subscribe{
//			Subscriptions: map[string]paho.SubscribeOptions{
//				"/in/ocpp2.0.1/cs1": {},
//			},
//		})
//		if err != nil {
//			t.Fatalf("subscribing to topic: %v", err)
//		}
//
//		cs := &registry.ChargeStation{
//			ClientId:             "cs1",
//			SecurityProfile:      registry.UnsecuredTransportWithBasicAuth,
//			Base64SHA256Password: "XohImNooBHFR0OVvjcYpJ3NgPQ1qq73WKhHvch0VQtg=", // password
//		}
//
//		mockRegistry := registry.NewMockRegistry()
//		mockRegistry.ChargeStations["cs1"] = cs
//
//		srv := httptest.NewServer(server.NewWebsocketHandler(
//			server.WithMqttBrokerUrl(addr),
//			server.WithDeviceRegistry(mockRegistry)))
//		defer srv.Close()
//
//		authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cs.ClientId, "password")))
//		dialOptions := &websocket.DialOptions{
//			Subprotocols: []string{"ocpp1.6", "ocpp2.0.1"},
//			HTTPHeader: http.Header{
//				"authorization": []string{authHeader},
//			},
//		}
//
//		conn, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/ws/cs1", srv.URL), dialOptions)
//		if err != nil {
//			t.Fatalf("dialing server: %v", err)
//		}
//		defer func() {
//			err := conn.Close(websocket.StatusNormalClosure, "OK")
//			if err != nil {
//				t.Logf("WARN: websocket close: %v", err)
//			}
//		}()
//
//		if conn.Subprotocol() != "ocpp2.0.1" {
//			t.Errorf("subprotocol: want %s, got %s", "ocpp2.0.1", conn.Subprotocol())
//		}
//
//		call := ocpp.Call{
//			MessageId: "1",
//			Action:    "EchoRequest",
//			Payload:   json.RawMessage("\"Payload\""),
//		}
//		data, err := json.Marshal(call)
//		if err != nil {
//			t.Fatalf("encoding call: %v", err)
//		}
//
//		err = conn.Write(ctx, websocket.MessageText, data)
//		if err != nil {
//			t.Fatalf("sending message: %v", err)
//		}
//
//		_ = conn.Close(websocket.StatusNormalClosure, "Something went wrong")
//	}
func createTestKeyPairs(t *testing.T, clientId string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, *ecdsa.PrivateKey) {
	caKeyPair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating CA key pair: %v", err)
	}

	caSerialNumber, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		t.Fatalf("generating CA serial number: %v", err)
	}

	caCertTemplate := x509.Certificate{
		SerialNumber: caSerialNumber,
		Subject: pkix.Name{
			CommonName:   "Test CA",
			Organization: []string{"Thoughtworks"},
		},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Minute),
	}

	caCertBytes, err :=
		x509.CreateCertificate(rand.Reader, &caCertTemplate, &caCertTemplate, &caKeyPair.PublicKey, caKeyPair)
	if err != nil {
		t.Fatalf("creating CA certificate: %v", err)
	}
	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		t.Fatalf("parsing CA certificate: %v", err)
	}

	clientKeyPair, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating client key pair: %v", err)
	}

	clientSerialNumber, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		t.Fatalf("generating client serial number: %v", err)
	}

	clientCertTemplate := x509.Certificate{
		SerialNumber: clientSerialNumber,
		Subject: pkix.Name{
			CommonName:   clientId,
			Organization: []string{"Thoughtworks"},
		},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Minute),
	}

	clientCertBytes, err := x509.CreateCertificate(rand.Reader, &clientCertTemplate, caCert, &clientKeyPair.PublicKey, caKeyPair)
	if err != nil {
		t.Fatalf("creating client certificate: %v", err)
	}
	clientCert, err := x509.ParseCertificate(clientCertBytes)
	if err != nil {
		t.Fatalf("parsing client certificate: %v", err)
	}

	return caCert, caKeyPair, clientCert, clientKeyPair
}
