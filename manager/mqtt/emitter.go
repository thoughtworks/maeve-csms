// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type Emitter interface {
	Emit(ctx context.Context, chargeStationId string, message *Message) error
}

type EmitterFunc func(ctx context.Context, chargeStationId string, message *Message) error

func (e EmitterFunc) Emit(ctx context.Context, chargeStationId string, message *Message) error {
	return e(ctx, chargeStationId, message)
}

type ProxyEmitter struct {
	emitter Emitter
}

func (p *ProxyEmitter) Emit(ctx context.Context, chargeStationId string, message *Message) error {
	return p.emitter.Emit(ctx, chargeStationId, message)
}

type MqttEmitter struct {
	conn        *autopaho.ConnectionManager
	mqttPrefix  string
	ocppVersion string
}

func (m *MqttEmitter) Emit(ctx context.Context, chargeStationId string, message *Message) error {
	topic := fmt.Sprintf("%s/out/%s/%s", m.mqttPrefix, m.ocppVersion, chargeStationId)
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshalling response of type %s: %v", message.Action, err)
	}
	_, err = m.conn.Publish(ctx, &paho.Publish{
		Topic:   topic,
		Payload: payload,
	})
	if err != nil {
		return fmt.Errorf("publishing to %s: %v", topic, err)
	}
	return nil
}

func NewMqttEmitter(conn *autopaho.ConnectionManager, mqttPrefix, ocppVersion string) Emitter {
	return &MqttEmitter{
		conn:        conn,
		mqttPrefix:  mqttPrefix,
		ocppVersion: ocppVersion,
	}
}
