// SPDX-License-Identifier: Apache-2.0

package mqtt_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/mqtt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"reflect"
	"regexp"
	"testing"
)

func TestBasicCallMaker(t *testing.T) {
	var got *mqtt.Message

	callMaker := mqtt.BasicCallMaker{
		E: mqtt.EmitterFunc(func(ctx context.Context, chargeStationId string, message *mqtt.Message) error {
			got = message
			return nil
		}),
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
	assert.Equal(t, mqtt.MessageTypeCall, got.MessageType)
	assert.Regexp(t, uuidPattern, got.MessageId)
	assert.Equal(t, "CertificateSigned", got.Action)
	assert.JSONEq(t, `{"certificateType":"V2GCertificate","certificateChain":"pemData"}`, string(got.RequestPayload))
}

func TestBasicCallMakerWithUnknownMessageType(t *testing.T) {
	var got *mqtt.Message

	callMaker := mqtt.BasicCallMaker{
		E: mqtt.EmitterFunc(func(ctx context.Context, chargeStationId string, message *mqtt.Message) error {
			got = message
			return nil
		}),
	}

	err := callMaker.Send(context.Background(), "cs001", &ocpp201.CertificateSignedRequestJson{})
	assert.ErrorContains(t, err, "unknown request type")
	assert.Nil(t, got)
}

func TestDataTransferCallMaker(t *testing.T) {
	var got *mqtt.Message

	callMaker := mqtt.DataTransferCallMaker{
		E: mqtt.EmitterFunc(func(ctx context.Context, chargeStationId string, message *mqtt.Message) error {
			got = message
			return nil
		}),
		Actions: map[reflect.Type]mqtt.DataTransferAction{
			reflect.TypeOf(&ocpp201.CertificateSignedRequestJson{}): {
				VendorId:  "org.openchargealliance.iso15118pnc",
				MessageId: "CertificateSigned",
			},
		},
	}

	certType := ocpp201.CertificateSigningUseEnumTypeV2GCertificate
	certChain := "pemData"
	err := callMaker.Send(context.Background(), "cs001", &ocpp201.CertificateSignedRequestJson{
		CertificateType:  &certType,
		CertificateChain: certChain,
	})
	require.NoError(t, err)

	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)
	assert.Equal(t, mqtt.MessageTypeCall, got.MessageType)
	assert.Regexp(t, uuidPattern, got.MessageId)
	assert.Equal(t, "DataTransfer", got.Action)
	assert.JSONEq(t, `{"vendorId":"org.openchargealliance.iso15118pnc","messageId":"CertificateSigned","data":"{\"certificateChain\":\"pemData\",\"certificateType\":\"V2GCertificate\"}"}`, string(got.RequestPayload))
}

func TestDataTransferCallMakerWithUnknownMessageType(t *testing.T) {
	var got *mqtt.Message

	callMaker := mqtt.DataTransferCallMaker{
		E: mqtt.EmitterFunc(func(ctx context.Context, chargeStationId string, message *mqtt.Message) error {
			got = message
			return nil
		}),
	}

	err := callMaker.Send(context.Background(), "cs001", &ocpp201.SignCertificateRequestJson{})
	assert.ErrorContains(t, err, "unknown request type")
	assert.Nil(t, got)
}

type cantMarshalType int

func (cantMarshalType) MarshalJSON() ([]byte, error) {
	return nil, errors.New("expected")
}

func (cantMarshalType) IsRequest() {}

func TestDataTransferCallMakerCantMarshalRequest(t *testing.T) {
	var got *mqtt.Message

	callMaker := mqtt.DataTransferCallMaker{
		E: mqtt.EmitterFunc(func(ctx context.Context, chargeStationId string, message *mqtt.Message) error {
			got = message
			return nil
		}),
		Actions: map[reflect.Type]mqtt.DataTransferAction{
			reflect.TypeOf(cantMarshalType(0)): {
				VendorId:  "CantMarshal",
				MessageId: "Me",
			},
		},
	}

	err := callMaker.Send(context.Background(), "cs001", cantMarshalType(1))
	assert.ErrorContains(t, err, "marshaling request")
	assert.Nil(t, got)
}
