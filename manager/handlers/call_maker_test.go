// SPDX-License-Identifier: Apache-2.0

package handlers_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"reflect"
	"regexp"
	"testing"
)

type FakeEmitter struct {
	called          bool
	ocppVersion     transport.OcppVersion
	chargeStationId string
	msg             *transport.Message
}

func (e *FakeEmitter) Emit(_ context.Context, ocppVersion transport.OcppVersion, chargeStationId string, message *transport.Message) error {
	e.called = true
	e.ocppVersion = ocppVersion
	e.chargeStationId = chargeStationId
	e.msg = message
	return nil
}

func TestCallMaker(t *testing.T) {
	emitter := &FakeEmitter{}
	callMaker := &handlers.OcppCallMaker{
		Emitter:     emitter,
		OcppVersion: transport.OcppVersion201,
		Actions: map[reflect.Type]string{
			reflect.TypeOf(&ocpp201.CertificateSignedRequestJson{}): "CertificateSigned",
		},
	}

	certType := ocpp201.CertificateSigningUseEnumTypeV2GCertificate
	certChain := "pemData"
	err := callMaker.Send(context.Background(), "cs001", &ocpp201.CertificateSignedRequestJson{
		CertificateType:  &certType,
		CertificateChain: certChain,
	})
	assert.NoError(t, err)

	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
	assert.Equal(t, transport.MessageTypeCall, emitter.msg.MessageType)
	assert.Regexp(t, uuidPattern, emitter.msg.MessageId)
	assert.Equal(t, "CertificateSigned", emitter.msg.Action)
	assert.JSONEq(t, `{"certificateType":"V2GCertificate","certificateChain":"pemData"}`, string(emitter.msg.RequestPayload))
}

func TestCallMakerWithUnknownMessageType(t *testing.T) {
	emitter := &FakeEmitter{}
	callMaker := &handlers.OcppCallMaker{
		Emitter:     emitter,
		OcppVersion: transport.OcppVersion201,
		Actions: map[reflect.Type]string{
			reflect.TypeOf(&ocpp201.CertificateSignedRequestJson{}): "CertificateSigned",
		},
	}

	err := callMaker.Send(context.Background(), "cs001", &ocpp201.AuthorizeRequestJson{})
	assert.ErrorContains(t, err, "unknown request type")
	assert.Nil(t, emitter.msg)
}
