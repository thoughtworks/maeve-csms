package ocpp16_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	handlers16 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"testing"
)

func TestDataTransferHandlerRoutesCall(t *testing.T) {
	dth := handlers16.DataTransferHandler{
		SchemaFS: schemas.OcppSchemas,
		CallRoutes: map[string]map[string]handlers.CallRoute{
			"org.openchargealliance.iso15118pnc": {
				"GetCertificateStatus": {
					NewRequest:     func() ocpp.Request { return new(types.GetCertificateStatusRequestJson) },
					RequestSchema:  "ocpp201/GetCertificateStatusRequest.json",
					ResponseSchema: "ocpp201/GetCertificateStatusResponse.json",
					Handler: handlers.CallHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
						ocspResult := "ocsp-result"
						return &types.GetCertificateStatusResponseJson{
							Status:     types.GetCertificateStatusEnumTypeAccepted,
							OcspResult: &ocspResult,
						}, nil
					}),
				},
			},
		},
	}

	getCertificateStatusMessageId := "GetCertificateStatus"
	getCertificateStatusData := "{\"ocspRequestData\":{\"hashAlgorithm\":\"SHA256\",\"issuerKeyHash\":\"key-hash\",\"issuerNameHash\":\"name-hash\",\"responderURL\":\"http://ocsp-qa.example.com:8080\",\"serialNumber\":\"serial-number\"}}"
	req := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &getCertificateStatusMessageId,
		Data:      &getCertificateStatusData,
	}

	got, err := dth.HandleCall(context.Background(), "cs001", req)
	require.NoError(t, err)

	expectedData := "{\"ocspResult\":\"ocsp-result\",\"status\":\"Accepted\"}"
	want := &ocpp16.DataTransferResponseJson{
		Data:   &expectedData,
		Status: ocpp16.DataTransferResponseJsonStatusAccepted,
	}

	assert.Equal(t, want, got)
}

func TestDataTransferHandlerWithUnknownVendorId(t *testing.T) {
	dth := handlers16.DataTransferHandler{
		CallRoutes: map[string]map[string]handlers.CallRoute{},
	}

	getCertificateStatusMessageId := "GetCertificateStatus"
	getCertificateStatusData := "{\"ocspRequestData\":{\"hashAlgorithm\":\"SHA256\",\"issuerKeyHash\":\"key-hash\",\"issuerNameHash\":\"name-hash\",\"responderURL\":\"http://ocsp-qa.example.com:8080\",\"serialNumber\":\"serial-number\"}}"
	req := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &getCertificateStatusMessageId,
		Data:      &getCertificateStatusData,
	}

	got, err := dth.HandleCall(context.Background(), "cs001", req)
	require.NoError(t, err)

	want := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusUnknownVendorId,
	}

	assert.Equal(t, want, got)
}

func TestDataTransferHandlerWithUnknownMessageId(t *testing.T) {
	dth := handlers16.DataTransferHandler{
		CallRoutes: map[string]map[string]handlers.CallRoute{
			"org.openchargealliance.iso15118pnc": {},
		},
	}

	getCertificateStatusMessageId := "GetCertificateStatus"
	getCertificateStatusData := "{\"ocspRequestData\":{\"hashAlgorithm\":\"SHA256\",\"issuerKeyHash\":\"key-hash\",\"issuerNameHash\":\"name-hash\",\"responderURL\":\"http://ocsp-qa.example.com:8080\",\"serialNumber\":\"serial-number\"}}"
	req := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &getCertificateStatusMessageId,
		Data:      &getCertificateStatusData,
	}

	got, err := dth.HandleCall(context.Background(), "cs001", req)
	require.NoError(t, err)

	want := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusUnknownMessageId,
	}

	assert.Equal(t, want, got)
}

func TestDataTransferHandlerWithEmptyResult(t *testing.T) {
	dth := handlers16.DataTransferHandler{
		SchemaFS: schemas.OcppSchemas,
		CallRoutes: map[string]map[string]handlers.CallRoute{
			"org.openchargealliance.iso15118pnc": {
				"GetCertificateStatus": {
					NewRequest:     func() ocpp.Request { return new(types.GetCertificateStatusRequestJson) },
					RequestSchema:  "ocpp201/GetCertificateStatusRequest.json",
					ResponseSchema: "ocpp201/GetCertificateStatusResponse.json",
					Handler: handlers.CallHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
						return nil, nil
					}),
				},
			},
		},
	}

	getCertificateStatusMessageId := "GetCertificateStatus"
	getCertificateStatusData := "{\"ocspRequestData\":{\"hashAlgorithm\":\"SHA256\",\"issuerKeyHash\":\"key-hash\",\"issuerNameHash\":\"name-hash\",\"responderURL\":\"http://ocsp-qa.example.com:8080\",\"serialNumber\":\"serial-number\"}}"
	req := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &getCertificateStatusMessageId,
		Data:      &getCertificateStatusData,
	}

	got, err := dth.HandleCall(context.Background(), "cs001", req)
	require.NoError(t, err)

	want := &ocpp16.DataTransferResponseJson{
		Status: ocpp16.DataTransferResponseJsonStatusAccepted,
	}

	assert.Equal(t, want, got)
}

func TestDataTransferHandlerWithErrorResult(t *testing.T) {
	dth := handlers16.DataTransferHandler{
		SchemaFS: schemas.OcppSchemas,
		CallRoutes: map[string]map[string]handlers.CallRoute{
			"org.openchargealliance.iso15118pnc": {
				"GetCertificateStatus": {
					NewRequest:     func() ocpp.Request { return new(types.GetCertificateStatusRequestJson) },
					RequestSchema:  "ocpp201/GetCertificateStatusRequest.json",
					ResponseSchema: "ocpp201/GetCertificateStatusResponse.json",
					Handler: handlers.CallHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
						return nil, errors.New("expected error")
					}),
				},
			},
		},
	}

	getCertificateStatusMessageId := "GetCertificateStatus"
	getCertificateStatusData := "{\"ocspRequestData\":{\"hashAlgorithm\":\"SHA256\",\"issuerKeyHash\":\"key-hash\",\"issuerNameHash\":\"name-hash\",\"responderURL\":\"http://ocsp-qa.example.com:8080\",\"serialNumber\":\"serial-number\"}}"
	req := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &getCertificateStatusMessageId,
		Data:      &getCertificateStatusData,
	}

	_, err := dth.HandleCall(context.Background(), "cs001", req)
	require.ErrorContains(t, err, "expected error")
}

func TestDataTransferHandlerCantUnmarshalRequest(t *testing.T) {
	dth := handlers16.DataTransferHandler{
		SchemaFS: schemas.OcppSchemas,
		CallRoutes: map[string]map[string]handlers.CallRoute{
			"org.openchargealliance.iso15118pnc": {
				"GetCertificateStatus": {
					NewRequest:     func() ocpp.Request { return new(noUnmarshalRequest) },
					RequestSchema:  "ocpp201/GetCertificateStatusRequest.json",
					ResponseSchema: "ocpp201/GetCertificateStatusResponse.json",
					Handler: handlers.CallHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
						ocspResult := "ocsp-result"
						return &types.GetCertificateStatusResponseJson{
							Status:     types.GetCertificateStatusEnumTypeAccepted,
							OcspResult: &ocspResult,
						}, nil
					}),
				},
			},
		},
	}

	getCertificateStatusMessageId := "GetCertificateStatus"
	getCertificateStatusData := "{\"ocspRequestData\":{\"hashAlgorithm\":\"SHA256\",\"issuerKeyHash\":\"key-hash\",\"issuerNameHash\":\"name-hash\",\"responderURL\":\"http://ocsp-qa.example.com:8080\",\"serialNumber\":\"serial-number\"}}"
	req := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &getCertificateStatusMessageId,
		Data:      &getCertificateStatusData,
	}

	_, err := dth.HandleCall(context.Background(), "cs001", req)
	require.ErrorContains(t, err, "unmarshalling org.openchargealliance.iso15118pnc:GetCertificateStatus data transfer data")
}

type ArrayResponse []string

func (*ArrayResponse) IsResponse() {}

func TestDataTransferHandlerCantMarshalResponse(t *testing.T) {
	dth := handlers16.DataTransferHandler{
		SchemaFS: schemas.OcppSchemas,
		CallRoutes: map[string]map[string]handlers.CallRoute{
			"org.openchargealliance.iso15118pnc": {
				"GetCertificateStatus": {
					NewRequest:     func() ocpp.Request { return new(types.GetCertificateStatusRequestJson) },
					RequestSchema:  "ocpp201/GetCertificateStatusRequest.json",
					ResponseSchema: "ocpp201/GetCertificateStatusResponse.json",
					Handler: handlers.CallHandlerFunc(func(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
						return new(noMarshalResponse), nil
					}),
				},
			},
		},
	}

	getCertificateStatusMessageId := "GetCertificateStatus"
	getCertificateStatusData := "{\"ocspRequestData\":{\"hashAlgorithm\":\"SHA256\",\"issuerKeyHash\":\"key-hash\",\"issuerNameHash\":\"name-hash\",\"responderURL\":\"http://ocsp-qa.example.com:8080\",\"serialNumber\":\"serial-number\"}}"
	req := &ocpp16.DataTransferJson{
		VendorId:  "org.openchargealliance.iso15118pnc",
		MessageId: &getCertificateStatusMessageId,
		Data:      &getCertificateStatusData,
	}

	_, err := dth.HandleCall(context.Background(), "cs001", req)
	require.ErrorContains(t, err, "marshalling org.openchargealliance.iso15118pnc:GetCertificateStatus data transfer data")
}

type noUnmarshalRequest struct{}

func (*noUnmarshalRequest) IsRequest() {}

func (*noUnmarshalRequest) UnmarshalJSON(data []byte) error {
	return errors.New("expected to fail")
}

type noUnmarshalResponse struct{}

func (*noUnmarshalResponse) IsResponse() {}

func (*noUnmarshalResponse) UnmarshalJSON(data []byte) error {
	return errors.New("expected to fail")
}

type noMarshalResponse struct{}

func (*noMarshalResponse) IsResponse() {}

func (*noMarshalResponse) MarshalJSON() ([]byte, error) {
	return nil, errors.New("expected to fail")
}
