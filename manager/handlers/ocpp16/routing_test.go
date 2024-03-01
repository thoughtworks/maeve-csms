// SPDX-License-Identifier: Apache-2.0

package ocpp16_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"regexp"
	"testing"
)

type FakeEmitter struct {
	got *transport.Message
}

func (e *FakeEmitter) Emit(ctx context.Context, ocppVersion transport.OcppVersion, chargeStationId string, message *transport.Message) error {
	e.got = message
	return nil
}

func TestDataTransferCallMaker(t *testing.T) {
	emitter := &FakeEmitter{}
	callMaker := ocpp16.NewDataTransferCallMaker(emitter)

	certType := ocpp201.CertificateSigningUseEnumTypeV2GCertificate
	certChain := "pemData"
	err := callMaker.Send(context.Background(), "cs001", &ocpp201.CertificateSignedRequestJson{
		CertificateType:  &certType,
		CertificateChain: certChain,
	})
	require.NoError(t, err)

	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
	assert.Equal(t, transport.MessageTypeCall, emitter.got.MessageType)
	assert.Regexp(t, uuidPattern, emitter.got.MessageId)
	assert.Equal(t, "DataTransfer", emitter.got.Action)
	assert.JSONEq(t, `{"vendorId":"org.openchargealliance.iso15118pnc","messageId":"CertificateSigned","data":"{\"certificateChain\":\"pemData\",\"certificateType\":\"V2GCertificate\"}"}`, string(emitter.got.RequestPayload))
}

func TestDataTransferCallMakerWithUnknownMessageType(t *testing.T) {
	emitter := &FakeEmitter{}
	callMaker := ocpp16.NewDataTransferCallMaker(emitter)

	err := callMaker.Send(context.Background(), "cs001", &ocpp201.AuthorizeRequestJson{})
	assert.ErrorContains(t, err, "unknown request type")
	assert.Nil(t, emitter.got)
}
