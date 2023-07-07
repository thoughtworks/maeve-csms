// SPDX-License-Identifier: Apache-2.0

package ocpp16_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	handlers16 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"testing"
)

func TestDataTransferResultHandler(t *testing.T) {
	var callHandled bool

	dtrh := handlers16.DataTransferResultHandler{
		SchemaFS: schemas.OcppSchemas,
		CallResultRoutes: map[string]map[string]handlers.CallResultRoute{
			"org.openchargealliance.iso15118pnc": {
				"CertificateSigned": {
					NewRequest:     func() ocpp.Request { return new(ocpp201.CertificateSignedRequestJson) },
					NewResponse:    func() ocpp.Response { return new(ocpp201.CertificateSignedResponseJson) },
					RequestSchema:  "ocpp201/CertificateSignedRequest.json",
					ResponseSchema: "ocpp201/CertificateSignedResponse.json",
					Handler: handlers.CallResultHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
						callHandled = true
						return nil
					}),
				},
			},
		},
	}

	messageId := "CertificateSigned"
	dataTransferData := "{\"certificateChain\":\"pemData\",\"certificateType\":\"V2GCertificate\"}"
	dataTransferRequest := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &messageId,
		Data:      &dataTransferData,
	}

	dataTransferResultData := "{\"status\":\"Accepted\"}"
	dataTransferResult := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusAccepted,
		Data:   &dataTransferResultData,
	}

	err := dtrh.HandleCallResult(context.Background(), "cs001", dataTransferRequest, dataTransferResult, "state")
	require.NoError(t, err)

	assert.True(t, callHandled)
}

func TestDataTransferResultHandlerErrorWithUnknownVendorId(t *testing.T) {
	dtrh := handlers16.DataTransferResultHandler{
		CallResultRoutes: map[string]map[string]handlers.CallResultRoute{},
	}

	messageId := "CertificateSigned"
	dataTransferData := "{\"certificateChain\":\"pemData\",\"certificateType\":\"V2GCertificate\"}"
	dataTransferRequest := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &messageId,
		Data:      &dataTransferData,
	}

	dataTransferResultData := "{\"status\":\"Accepted\"}"
	dataTransferResult := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusAccepted,
		Data:   &dataTransferResultData,
	}

	err := dtrh.HandleCallResult(context.Background(), "cs001", dataTransferRequest, dataTransferResult, "state")
	require.ErrorContains(t, err, "unknown data transfer result vendor")
}

func TestDataTransferResultHandlerErrorWithUnknownMessageId(t *testing.T) {
	dtrh := handlers16.DataTransferResultHandler{
		CallResultRoutes: map[string]map[string]handlers.CallResultRoute{
			"org.openchargealliance.iso15118pnc": {},
		},
	}

	messageId := "CertificateSigned"
	dataTransferData := "{\"certificateChain\":\"pemData\",\"certificateType\":\"V2GCertificate\"}"
	dataTransferRequest := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &messageId,
		Data:      &dataTransferData,
	}

	dataTransferResultData := "{\"status\":\"Accepted\"}"
	dataTransferResult := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusAccepted,
		Data:   &dataTransferResultData,
	}

	err := dtrh.HandleCallResult(context.Background(), "cs001", dataTransferRequest, dataTransferResult, "state")
	require.ErrorContains(t, err, "unknown data transfer result message id")
}

func TestDataTransferResultHandlerErrorWhenRequestPayloadIsInvalid(t *testing.T) {
	dtrh := handlers16.DataTransferResultHandler{
		SchemaFS: schemas.OcppSchemas,
		CallResultRoutes: map[string]map[string]handlers.CallResultRoute{
			"org.openchargealliance.iso15118pnc": {
				"CertificateSigned": {
					NewRequest:     func() ocpp.Request { return new(ocpp201.CertificateSignedRequestJson) },
					NewResponse:    func() ocpp.Response { return new(ocpp201.CertificateSignedResponseJson) },
					RequestSchema:  "ocpp201/CertificateSignedRequest.json",
					ResponseSchema: "ocpp201/CertificateSignedResponse.json",
					Handler: handlers.CallResultHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
						return nil
					}),
				},
			},
		},
	}

	messageId := "CertificateSigned"
	dataTransferData := "{\"certificateType\":\"V2GCertificate\"}"
	dataTransferRequest := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &messageId,
		Data:      &dataTransferData,
	}

	dataTransferResultData := "{\"status\":\"Accepted\"}"
	dataTransferResult := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusAccepted,
		Data:   &dataTransferResultData,
	}

	err := dtrh.HandleCallResult(context.Background(), "cs001", dataTransferRequest, dataTransferResult, "state")
	require.ErrorContains(t, err, "validating org.openchargealliance.iso15118pnc:CertificateSigned data transfer result request data")
}

func TestDataTransferResultHandlerErrorWhenResponsePayloadIsInvalid(t *testing.T) {
	dtrh := handlers16.DataTransferResultHandler{
		SchemaFS: schemas.OcppSchemas,
		CallResultRoutes: map[string]map[string]handlers.CallResultRoute{
			"org.openchargealliance.iso15118pnc": {
				"CertificateSigned": {
					NewRequest:     func() ocpp.Request { return new(ocpp201.CertificateSignedRequestJson) },
					NewResponse:    func() ocpp.Response { return new(ocpp201.CertificateSignedResponseJson) },
					RequestSchema:  "ocpp201/CertificateSignedRequest.json",
					ResponseSchema: "ocpp201/CertificateSignedResponse.json",
					Handler: handlers.CallResultHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
						return nil
					}),
				},
			},
		},
	}

	messageId := "CertificateSigned"
	dataTransferData := "{\"certificateChain\":\"pemData\",\"certificateType\":\"V2GCertificate\"}"
	dataTransferRequest := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &messageId,
		Data:      &dataTransferData,
	}

	dataTransferResultData := "{}"
	dataTransferResult := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusAccepted,
		Data:   &dataTransferResultData,
	}

	err := dtrh.HandleCallResult(context.Background(), "cs001", dataTransferRequest, dataTransferResult, "state")
	require.ErrorContains(t, err, "validating org.openchargealliance.iso15118pnc:CertificateSigned data transfer result response data")
}

func TestDataTransferResultHandlerErrorWhenCantUnmarshalRequest(t *testing.T) {
	dtrh := handlers16.DataTransferResultHandler{
		SchemaFS: schemas.OcppSchemas,
		CallResultRoutes: map[string]map[string]handlers.CallResultRoute{
			"org.openchargealliance.iso15118pnc": {
				"CertificateSigned": {
					NewRequest:     func() ocpp.Request { return new(noUnmarshalRequest) },
					NewResponse:    func() ocpp.Response { return new(ocpp201.CertificateSignedResponseJson) },
					RequestSchema:  "ocpp201/CertificateSignedRequest.json",
					ResponseSchema: "ocpp201/CertificateSignedResponse.json",
					Handler: handlers.CallResultHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
						return nil
					}),
				},
			},
		},
	}

	messageId := "CertificateSigned"
	dataTransferData := "{\"certificateChain\":\"pemData\",\"certificateType\":\"V2GCertificate\"}"
	dataTransferRequest := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &messageId,
		Data:      &dataTransferData,
	}

	dataTransferResultData := "{\"status\":\"Accepted\"}"
	dataTransferResult := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusAccepted,
		Data:   &dataTransferResultData,
	}

	err := dtrh.HandleCallResult(context.Background(), "cs001", dataTransferRequest, dataTransferResult, "state")
	require.ErrorContains(t, err, "unmarshalling org.openchargealliance.iso15118pnc:CertificateSigned data transfer request data")
}

func TestDataTransferResultHandlerErrorWhenCantUnmarshalResponse(t *testing.T) {
	dtrh := handlers16.DataTransferResultHandler{
		SchemaFS: schemas.OcppSchemas,
		CallResultRoutes: map[string]map[string]handlers.CallResultRoute{
			"org.openchargealliance.iso15118pnc": {
				"CertificateSigned": {
					NewRequest:     func() ocpp.Request { return new(ocpp201.CertificateSignedRequestJson) },
					NewResponse:    func() ocpp.Response { return new(noUnmarshalResponse) },
					RequestSchema:  "ocpp201/CertificateSignedRequest.json",
					ResponseSchema: "ocpp201/CertificateSignedResponse.json",
					Handler: handlers.CallResultHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
						return nil
					}),
				},
			},
		},
	}

	messageId := "CertificateSigned"
	dataTransferData := "{\"certificateChain\":\"pemData\",\"certificateType\":\"V2GCertificate\"}"
	dataTransferRequest := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &messageId,
		Data:      &dataTransferData,
	}

	dataTransferResultData := "{\"status\":\"Accepted\"}"
	dataTransferResult := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusAccepted,
		Data:   &dataTransferResultData,
	}

	err := dtrh.HandleCallResult(context.Background(), "cs001", dataTransferRequest, dataTransferResult, "state")
	require.ErrorContains(t, err, "unmarshalling org.openchargealliance.iso15118pnc:CertificateSigned data transfer response data")
}
