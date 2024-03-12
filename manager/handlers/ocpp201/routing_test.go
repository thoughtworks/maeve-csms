// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	clockTest "k8s.io/utils/clock/testing"
	"testing"
	"time"
)

type fakeEmitter struct{}

func (f fakeEmitter) Emit(ctx context.Context, ocppVersion transport.OcppVersion, chargeStationId string, message *transport.Message) error {
	return nil
}

type fakeTariffService struct{}

func (f fakeTariffService) CalculateCost(transaction *store.Transaction) (float64, error) {
	return 42.0, nil
}

type fakeCertValidationService struct{}

func (f fakeCertValidationService) ValidatePEMCertificateChain(ctx context.Context, pemChain []byte, eMAID string) (*string, error) {
	return nil, nil
}

func (f fakeCertValidationService) ValidateHashedCertificateChain(ctx context.Context, ocspRequestData []types.OCSPRequestDataType) (*string, error) {
	return nil, nil
}

type fakeChargeStationCertProvider struct{}

func (f fakeChargeStationCertProvider) ProvideCertificate(ctx context.Context, typ services.CertificateType, pemEncodedCSR string, csId string) (pemEncodedCertificateChain string, err error) {
	return "", nil
}

type fakeContractCertProvider struct{}

func (f fakeContractCertProvider) ProvideCertificate(ctx context.Context, exiRequest string) (services.EvCertificate15118Response, error) {
	return services.EvCertificate15118Response{
		Status:                     types.Iso15118EVCertificateStatusEnumTypeAccepted,
		CertificateInstallationRes: "",
	}, nil
}

func TestRoutingCalls(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)
	clock := clockTest.NewFakePassiveClock(now)

	engine := inmemory.NewStore(clock)

	router := ocpp201.NewRouter(&fakeEmitter{},
		clock,
		engine,
		&fakeTariffService{},
		&fakeCertValidationService{},
		&fakeChargeStationCertProvider{},
		&fakeContractCertProvider{},
		5*time.Minute,
		schemas.OcppSchemas,
	)

	inputMessages := map[string]ocpp.Request{
		"Authorize": &types.AuthorizeRequestJson{
			IdToken: types.IdTokenType{
				Type:    types.IdTokenEnumTypeISO14443,
				IdToken: "DEADBEEF",
			},
		},
		"BootNotification": &types.BootNotificationRequestJson{
			ChargingStation: types.ChargingStationType{
				Model:        "Powergen",
				SerialNumber: makePtr("012345ABCDEF"),
				VendorName:   "Vendor",
			},
			Reason: types.BootReasonEnumTypePowerUp,
		},
		"FirmwareStatusNotification": &types.FirmwareStatusNotificationRequestJson{
			Status: types.FirmwareStatusEnumTypeDownloading,
		},
		"Get15118EVCertificate": &types.Get15118EVCertificateRequestJson{
			Action:                types.CertificateActionEnumTypeInstall,
			ExiRequest:            "",
			Iso15118SchemaVersion: "15118-2",
		},
		"GetCertificateStatus": &types.GetCertificateStatusRequestJson{
			OcspRequestData: types.OCSPRequestDataType{
				HashAlgorithm:  "SHA256",
				IssuerKeyHash:  "123456ABCDEF",
				IssuerNameHash: "ABCDEF123456",
				ResponderURL:   "https://example.org",
				SerialNumber:   "123456789",
			},
		},
		"Heartbeat": &types.HeartbeatRequestJson{},
		"LogStatusNotification": &types.LogStatusNotificationRequestJson{
			Status: types.UploadLogStatusEnumTypeUploadFailure,
		},
		"MeterValues": &types.MeterValuesRequestJson{
			EvseId: 1,
			MeterValue: []types.MeterValueType{
				{
					SampledValue: []types.SampledValueType{
						{
							Location:  makePtr(types.LocationEnumTypeOutlet),
							Measurand: makePtr(types.MeasurandEnumTypeCurrentExport),
							Value:     12,
						},
					},
					Timestamp: "2023-06-15T15:05:00+01:00",
				},
			},
		},
		"SecurityEventNotification": &types.SecurityEventNotificationRequestJson{
			Timestamp: "2023-06-15T15:05:00+01:00",
			Type:      "SettingSystemTime",
		},
		"SignCertificate": &types.SignCertificateRequestJson{
			CertificateType: makePtr(types.CertificateSigningUseEnumTypeChargingStationCertificate),
			Csr:             "",
		},
		"StatusNotification": &types.StatusNotificationRequestJson{
			ConnectorId:     1,
			ConnectorStatus: types.ConnectorStatusEnumTypeAvailable,
			EvseId:          1,
			Timestamp:       "2023-06-15T15:05:00+01:00",
		},
		"TransactionEvent": &types.TransactionEventRequestJson{
			EventType: types.TransactionEventEnumTypeStarted,
			Evse: &types.EVSEType{
				Id:          1,
				ConnectorId: makePtr(1),
			},
			IdToken: &types.IdTokenType{
				IdToken: "DEADBEEF",
				Type:    types.IdTokenEnumTypeISO14443,
			},
			MeterValue: []types.MeterValueType{
				{
					Timestamp: "2023-06-15T15:05:00+01:00",
					SampledValue: []types.SampledValueType{
						{
							Location:  makePtr(types.LocationEnumTypeOutlet),
							Measurand: makePtr(types.MeasurandEnumTypeCurrentExport),
							Value:     24,
						},
					},
				},
			},
			NumberOfPhasesUsed: makePtr(3),
			Offline:            false,
			SeqNo:              1,
			Timestamp:          "2023-06-15T15:05:00+01:00",
			TransactionInfo:    types.TransactionType{},
			TriggerReason:      types.TriggerReasonEnumTypeAuthorized,
		},
	}

	for action, req := range inputMessages {
		t.Run(action, func(t *testing.T) {
			reqBytes, err := json.Marshal(req)
			require.NoError(t, err)

			messageId, err := uuid.NewUUID()
			require.NoError(t, err)

			msg := transport.Message{
				MessageType:    transport.MessageTypeCall,
				Action:         action,
				MessageId:      messageId.String(),
				RequestPayload: reqBytes,
			}
			err = router.Route(context.TODO(), "cs001", msg)
			assert.NoError(t, err)
		})
	}
}

func TestRoutingCallResults(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)
	clock := clockTest.NewFakePassiveClock(now)

	engine := inmemory.NewStore(clock)

	router := ocpp201.NewRouter(&fakeEmitter{},
		clock,
		engine,
		&fakeTariffService{},
		&fakeCertValidationService{},
		&fakeChargeStationCertProvider{},
		&fakeContractCertProvider{},
		5*time.Minute,
		schemas.OcppSchemas,
	)

	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("test"),
	}
	pemBytes := pem.EncodeToMemory(pemBlock)

	inputMessages := map[string]struct {
		request  ocpp.Request
		response ocpp.Response
	}{
		"CertificateSigned": {
			request: &types.CertificateSignedRequestJson{
				CertificateChain: string(pemBytes),
				CertificateType:  makePtr(types.CertificateSigningUseEnumTypeV2GCertificate),
			},
			response: &types.CertificateSignedResponseJson{
				Status: types.CertificateSignedStatusEnumTypeRejected,
			},
		},
		"InstallCertificate": {
			request: &types.InstallCertificateRequestJson{
				Certificate:     string(pemBytes),
				CertificateType: types.InstallCertificateUseEnumTypeMORootCertificate,
			},
			response: &types.InstallCertificateResponseJson{
				Status: types.InstallCertificateStatusEnumTypeAccepted,
			},
		},
		"SetVariables": {
			request: &types.SetVariablesRequestJson{
				SetVariableData: []types.SetVariableDataType{
					{
						Component: types.ComponentType{
							Name: "AlignedDataCtrlr",
						},
						Variable: types.VariableType{
							Name: "Interval",
						},
						AttributeValue: "60",
					},
				},
			},
			response: &types.SetVariablesResponseJson{
				SetVariableResult: []types.SetVariableResultType{
					{
						Component: types.ComponentType{
							Name: "AlignedDataCtrlr",
						},
						Variable: types.VariableType{
							Name: "Interval",
						},
						AttributeStatus: types.SetVariableStatusEnumTypeAccepted,
					},
				},
			},
		},
		"TriggerMessage": {
			request: &types.TriggerMessageRequestJson{
				RequestedMessage: types.MessageTriggerEnumTypeHeartbeat,
			},
			response: &types.TriggerMessageResponseJson{
				Status: types.TriggerMessageStatusEnumTypeAccepted,
			},
		},
	}

	for action, input := range inputMessages {
		t.Run(action, func(t *testing.T) {
			reqBytes, err := json.Marshal(input.request)
			require.NoError(t, err)

			respBytes, err := json.Marshal(input.response)
			require.NoError(t, err)

			messageId, err := uuid.NewUUID()
			require.NoError(t, err)

			msg := transport.Message{
				MessageType:     transport.MessageTypeCallResult,
				Action:          action,
				MessageId:       messageId.String(),
				RequestPayload:  reqBytes,
				ResponsePayload: respBytes,
			}

			err = router.Route(context.TODO(), "cs001", msg)
			assert.NoError(t, err)
		})
	}
}
