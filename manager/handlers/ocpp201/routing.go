// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"io/fs"
	"k8s.io/utils/clock"
	"reflect"
	"time"
)

func NewRouter(emitter transport.Emitter,
	clk clock.PassiveClock,
	engine store.Engine,
	tariffService services.TariffService,
	certValidationService services.CertificateValidationService,
	chargeStationCertProvider services.ChargeStationCertificateProvider,
	contractCertProvider services.ContractCertificateProvider,
	heartbeatInterval time.Duration,
	schemaFS fs.FS) transport.MessageHandler {

	return &handlers.Router{
		Emitter:     emitter,
		SchemaFS:    schemaFS,
		OcppVersion: transport.OcppVersion201,
		CallRoutes: map[string]handlers.CallRoute{
			"Authorize": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.AuthorizeRequestJson) },
				RequestSchema:  "ocpp201/AuthorizeRequest.json",
				ResponseSchema: "ocpp201/AuthorizeResponse.json",
				Handler: AuthorizeHandler{
					TokenStore:                   engine,
					CertificateValidationService: certValidationService,
				},
			},
			"BootNotification": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.BootNotificationRequestJson) },
				RequestSchema:  "ocpp201/BootNotificationRequest.json",
				ResponseSchema: "ocpp201/BootNotificationResponse.json",
				Handler: BootNotificationHandler{
					Clock:               clk,
					HeartbeatInterval:   int(heartbeatInterval.Seconds()),
					RuntimeDetailsStore: engine,
				},
			},
			"FirmwareStatusNotification": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.FirmwareStatusNotificationRequestJson) },
				RequestSchema:  "ocpp201/FirmwareStatusNotificationRequest.json",
				ResponseSchema: "ocpp201/FirmwareStatusNotificationResponse.json",
				Handler:        FirmwareStatusNotificationHandler{},
			},
			"GetCertificateStatus": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.GetCertificateStatusRequestJson) },
				RequestSchema:  "ocpp201/GetCertificateStatusRequest.json",
				ResponseSchema: "ocpp201/GetCertificateStatusResponse.json",
				Handler: GetCertificateStatusHandler{
					CertificateValidationService: certValidationService,
				},
			},
			"Get15118EVCertificate": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.Get15118EVCertificateRequestJson) },
				RequestSchema:  "ocpp201/Get15118EVCertificateRequest.json",
				ResponseSchema: "ocpp201/Get15118EVCertificateResponse.json",
				Handler: Get15118EvCertificateHandler{
					ContractCertificateProvider: contractCertProvider,
				},
			},
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				RequestSchema:  "ocpp201/HeartbeatRequest.json",
				ResponseSchema: "ocpp201/HeartbeatResponse.json",
				Handler: HeartbeatHandler{
					Clock: clk,
				},
			},
			"LogStatusNotification": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.LogStatusNotificationRequestJson) },
				RequestSchema:  "ocpp201/LogStatusNotificationRequest.json",
				ResponseSchema: "ocpp201/LogStatusNotificationResponse.json",
				Handler:        LogStatusNotificationHandler{},
			},
			"MeterValues": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.MeterValuesRequestJson) },
				RequestSchema:  "ocpp201/MeterValuesRequest.json",
				ResponseSchema: "ocpp201/MeterValuesResponse.json",
				Handler:        MeterValuesHandler{},
			},
			"StatusNotification": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.StatusNotificationRequestJson) },
				RequestSchema:  "ocpp201/StatusNotificationRequest.json",
				ResponseSchema: "ocpp201/StatusNotificationResponse.json",
				Handler:        handlers.CallHandlerFunc(StatusNotificationHandler),
			},
			"SignCertificate": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.SignCertificateRequestJson) },
				RequestSchema:  "ocpp201/SignCertificateRequest.json",
				ResponseSchema: "ocpp201/SignCertificateResponse.json",
				Handler: SignCertificateHandler{
					ChargeStationCertificateProvider: chargeStationCertProvider,
					Store:                            engine,
				},
			},
			"SecurityEventNotification": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.SecurityEventNotificationRequestJson) },
				RequestSchema:  "ocpp201/SecurityEventNotificationRequest.json",
				ResponseSchema: "ocpp201/SecurityEventNotificationResponse.json",
				Handler:        SecurityEventNotificationHandler{},
			},
			"TransactionEvent": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.TransactionEventRequestJson) },
				RequestSchema:  "ocpp201/TransactionEventRequest.json",
				ResponseSchema: "ocpp201/TransactionEventResponse.json",
				Handler: TransactionEventHandler{
					Store:         engine,
					TariffService: tariffService,
				},
			},
		},
		CallResultRoutes: map[string]handlers.CallResultRoute{
			"CertificateSigned": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.CertificateSignedRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.CertificateSignedResponseJson) },
				RequestSchema:  "ocpp201/CertificateSignedRequest.json",
				ResponseSchema: "ocpp201/CertificateSignedResponse.json",
				Handler: CertificateSignedResultHandler{
					Store: engine,
				},
			},
			"ChangeAvailability": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.ChangeAvailabilityRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.ChangeAvailabilityResponseJson) },
				RequestSchema:  "ocpp201/ChangeAvailabilityRequest.json",
				ResponseSchema: "ocpp201/ChangeAvailabilityResponse.json",
				Handler:        ChangeAvailabilityResultHandler{},
			},
			"GetTransactionStatus": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.GetTransactionStatusRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.GetTransactionStatusResponseJson) },
				RequestSchema:  "ocpp201/GetTransactionStatusRequest.json",
				ResponseSchema: "ocpp201/GetTransactionStatusResponse.json",
				Handler:        GetTransactionStatusResultHandler{},
			},
			"GetVariables": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.GetVariablesRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.GetVariablesResponseJson) },
				RequestSchema:  "ocpp201/GetVariablesRequest.json",
				ResponseSchema: "ocpp201/GetVariablesResponse.json",
				Handler:        GetVariablesResultHandler{},
			},
			"InstallCertificate": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.InstallCertificateRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.InstallCertificateResponseJson) },
				RequestSchema:  "ocpp201/InstallCertificateRequest.json",
				ResponseSchema: "ocpp201/InstallCertificateResponse.json",
				Handler: InstallCertificateResultHandler{
					Store: engine,
				},
			},
			"RequestStartTransaction": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.RequestStartTransactionRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.RequestStartTransactionResponseJson) },
				RequestSchema:  "ocpp201/RequestStartTransactionRequest.json",
				ResponseSchema: "ocpp201/RequestStartTransactionResponse.json",
				Handler:        RequestStartTransactionResultHandler{},
			},
			"RequestStopTransaction": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.RequestStopTransactionRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.RequestStopTransactionResponseJson) },
				RequestSchema:  "ocpp201/RequestStopTransactionRequest.json",
				ResponseSchema: "ocpp201/RequestStopTransactionResponse.json",
				Handler:        RequestStopTransactionResultHandler{},
			},
			"SetVariables": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.SetVariablesRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.SetVariablesResponseJson) },
				RequestSchema:  "ocpp201/SetVariablesRequest.json",
				ResponseSchema: "ocpp201/SetVariablesResponse.json",
				Handler: SetVariablesResultHandler{
					Store: engine,
				},
			},
			"TriggerMessage": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.TriggerMessageRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.TriggerMessageResponseJson) },
				RequestSchema:  "ocpp201/TriggerMessageRequest.json",
				ResponseSchema: "ocpp201/TriggerMessageResponse.json",
				Handler: TriggerMessageResultHandler{
					Store: engine,
				},
			},
			"UnlockConnector": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.UnlockConnectorRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.UnlockConnectorResponseJson) },
				RequestSchema:  "ocpp201/UnlockConnectorRequest.json",
				ResponseSchema: "ocpp201/UnclockConnectorResponse.json",
				Handler:        UnlockConnectorResultHandler{},
			},
		},
	}
}

func NewCallMaker(e transport.Emitter) *handlers.OcppCallMaker {
	return &handlers.OcppCallMaker{
		Emitter:     e,
		OcppVersion: transport.OcppVersion201,
		Actions: map[reflect.Type]string{
			reflect.TypeOf(&ocpp201.CertificateSignedRequestJson{}):       "CertificateSigned",
			reflect.TypeOf(&ocpp201.ChangeAvailabilityRequestJson{}):      "ChangeAvailability",
			reflect.TypeOf(&ocpp201.GetTransactionStatusRequestJson{}):    "GetTransactionStatus",
			reflect.TypeOf(&ocpp201.GetVariablesRequestJson{}):            "GetVariables",
			reflect.TypeOf(&ocpp201.InstallCertificateRequestJson{}):      "InstallCertificate",
			reflect.TypeOf(&ocpp201.RequestStartTransactionRequestJson{}): "RequestStartTransaction",
			reflect.TypeOf(&ocpp201.RequestStopTransactionRequestJson{}):  "RequestStopTransaction",
			reflect.TypeOf(&ocpp201.SetVariablesRequestJson{}):            "SetVariables",
			reflect.TypeOf(&ocpp201.TriggerMessageRequestJson{}):          "TriggerMessage",
			reflect.TypeOf(&ocpp201.UnlockConnectorRequestJson{}):         "UnlockConnector",
		},
	}
}
