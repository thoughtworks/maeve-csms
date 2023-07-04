package server

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twlabs/ocpp2-broker-core/gateway/ocpp"
	"github.com/twlabs/ocpp2-broker-core/gateway/pipe"
	"github.com/twlabs/ocpp2-broker-core/gateway/registry"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"net/http"
	"net/url"
	"nhooyr.io/websocket"
	"time"
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

func NewWebsocketHandler(opts ...WebsocketOpt) http.Handler {
	s := new(WebsocketHandler)

	for _, opt := range opts {
		opt(s)
	}

	ensureDefaults(s)

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)
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
}

func (s *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientId := chi.URLParam(r, "id")
	if clientId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	cs := s.deviceRegistry.LookupChargeStation(clientId)
	if cs == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	switch cs.SecurityProfile {
	case registry.UnsecuredTransportWithBasicAuth:
		if r.TLS != nil || !checkAuthorization(r, cs) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	case registry.TLSWithBasicAuth:
		if r.TLS == nil || !checkAuthorization(r, cs) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	case registry.TLSWithClientSideCertificates:
		if r.TLS == nil || !checkCertificate(r, s.orgNames, cs) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{Subprotocols: []string{"ocpp2.0.1", "ocpp1.6"}})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	protocol := wsConn.Subprotocol()
	if protocol == "" {
		protocol = "ocpp2.0.1"
	}

	p := pipe.NewPipe(s.pipeOptions...)
	p.Start()
	defer p.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var mqttConn *autopaho.ConnectionManager
	mqttConn, err = autopaho.NewConnection(ctx, autopaho.ClientConfig{
		BrokerUrls:        s.mqttBrokerURLs,
		KeepAlive:         s.mqttKeepAliveInterval,
		ConnectRetryDelay: s.mqttConnectRetryDelay,
		//Debug:             LogLogger{},
		//PahoDebug:         LogLogger{},
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			log.Printf("connection up....")
			topicName := fmt.Sprintf("%s/out/%s/%s", s.mqttTopicPrefix, protocol, clientId)
			_, err = manager.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					topicName: {},
				},
			})
			if err != nil {
				log.Printf("subscribing to mqtt topic %s: %v", topicName, err)
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
					log.Printf("unmarshalling CSMS message: %v", err)
					return
				}
				p.CSMSRx <- &msg
			}),
			OnServerDisconnect: func(disconnect *paho.Disconnect) {
				log.Printf("server disconnect...")
			},
		},
	})
	if err != nil {
		log.Printf("connecting to mqtt on %s: %v", s.mqttBrokerURLs, err)
		_ = wsConn.Close(websocket.StatusProtocolError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer func() {
		err := mqttConn.Disconnect(context.Background())
		if err != nil {
			log.Printf("disconnecting from mqtt: %v", err)
		}
	}()
	err = mqttConn.AwaitConnection(ctx)
	if err != nil {
		log.Printf("waiting for mqtt on %s: %v", s.mqttBrokerURLs, err)
		return
	}

	// listen on the CSMS Tx channel and publish those messages on the inbound topic
	goPublishToCSMS(ctx, p.CSMSTx, mqttConn, s.mqttTopicPrefix, protocol, clientId)

	// listen the CS Tx channel and write those messages to the websocket
	goWriteToChargeStation(ctx, p.ChargeStationTx, wsConn)

	// read from the websocket and send to the CS Rx channel (CS Tx used for error)
	readFromChargeStation(ctx, wsConn, p.ChargeStationRx, p.ChargeStationTx)

	log.Printf("websocket handler complete")
}

func checkAuthorization(r *http.Request, cs *registry.ChargeStation) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}
	if username != cs.ClientId {
		return false
	}
	sha256pw := sha256.Sum256([]byte(password))
	b64sha256 := base64.StdEncoding.EncodeToString(sha256pw[:])
	return b64sha256 == cs.Base64SHA256Password
}

func checkCertificate(r *http.Request, orgNames []string, cs *registry.ChargeStation) bool {
	if len(r.TLS.PeerCertificates) == 0 {
		return false
	}

	leafCertificate := r.TLS.PeerCertificates[0]

	foundOrg := false
	for _, org := range leafCertificate.Subject.Organization {
		log.Printf("Org Name: %s; Allowed Org Names: %s", org, orgNames)

		if slices.Contains(orgNames, org) {
			foundOrg = true
			break
		}
	}
	if !foundOrg {
		return false
	}

	log.Printf("Client Id: %s", leafCertificate.Subject.CommonName)

	return cs.ClientId == leafCertificate.Subject.CommonName
}

func goPublishToCSMS(ctx context.Context, csmsTx chan *pipe.GatewayMessage, mqttConn *autopaho.ConnectionManager, topicPrefix, protocol, clientId string) {
	go func() {
		for {
			select {
			case msg := <-csmsTx:
				topic := fmt.Sprintf("%s/in/%s/%s", topicPrefix, protocol, clientId)
				//log.Printf("publishing message on %s: %v", topic, msg)
				data, err := json.Marshal(msg)
				if err != nil {
					log.Printf("marshaling message for publication: %v", err)
					continue
				}
				_, err = mqttConn.Publish(ctx, &paho.Publish{
					Topic:   topic,
					Payload: data,
					Properties: &paho.PublishProperties{
						ContentType:   "application/json",
						ResponseTopic: fmt.Sprintf("%s/out/%s/%s", topicPrefix, protocol, clientId),
					},
				})
				if err != nil {
					log.Printf("publishing message: %v\n", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func goWriteToChargeStation(ctx context.Context, chargeStationTx chan *pipe.GatewayMessage, wsConn *websocket.Conn) {
	go func() {
		for {
			select {
			case msg := <-chargeStationTx:
				data, err := marshalGatewayMessageAsOcpp(msg)
				if err != nil {
					log.Printf("marshaling gateway message for charge station: %v", err)
					continue
				}
				err = wsConn.Write(ctx, websocket.MessageText, data)
				if err != nil {
					log.Printf("writing to charge station: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func readFromChargeStation(ctx context.Context, wsConn *websocket.Conn, csRx chan *pipe.GatewayMessage, csTx chan *pipe.GatewayMessage) {
	for {
		typ, b, err := wsConn.Read(ctx)
		if errors.Is(err, context.DeadlineExceeded) {
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		} else if status := websocket.CloseStatus(err); status != -1 {
			log.Printf("connection closed with status %d", status)
			return
		} else if errors.Is(err, io.EOF) {
			log.Printf("connection closed")
			return
		} else if err != nil {
			log.Printf("reading from cs: %v", err)
			return
		}

		if typ != websocket.MessageText {
			msg := pipe.GatewayMessage{
				MessageType:      ocpp.MessageTypeCallError,
				MessageId:        "-1",
				ErrorCode:        ocpp.ErrorRpcFrameworkError,
				ErrorDescription: "websocket message type is not text",
			}
			csTx <- &msg
			continue
		}

		msg, err := unmarshalOcppAsGatewayMessage(b)
		if err != nil {
			msg := pipe.GatewayMessage{
				MessageType:      ocpp.MessageTypeCallError,
				MessageId:        "-1",
				ErrorCode:        ocpp.ErrorRpcFrameworkError,
				ErrorDescription: err.Error(),
			}
			csTx <- &msg
			continue
		}
		//log.Printf("received message: %v", msg)
		csRx <- msg
	}
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

type LogLogger struct{}

func (l LogLogger) Println(v ...interface{}) {
	log.Println(v...)
}

func (l LogLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
