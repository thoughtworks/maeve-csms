// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/has2be"
	"io/fs"
	"reflect"
	"time"

	"github.com/santhosh-tekuri/jsonschema"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	handlersHasToBe "github.com/thoughtworks/maeve-csms/manager/handlers/has2be"
	handlers16 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
)

type Router struct {
	CallRoutes       map[string]handlers.CallRoute
	CallResultRoutes map[string]handlers.CallResultRoute
}

func NewV16Router(emitter Emitter,
	clk clock.PassiveClock,
	engine store.Engine,
	certValidationService services.CertificateValidationService,
	chargeStationCertProvider services.ChargeStationCertificateProvider,
	contractCertProvider services.ContractCertificateProvider,
	heartbeatInterval time.Duration,
	schemaFS fs.FS) *Router {

	dataTransferCallMaker := DataTransferCallMaker{
		E: emitter,
		Actions: map[reflect.Type]DataTransferAction{
			reflect.TypeOf(&ocpp201.CertificateSignedRequestJson{}): {
				VendorId:  "org.openchargealliance.iso15118pnc",
				MessageId: "CertificateSigned",
			},
			reflect.TypeOf(&has2be.CertificateSignedRequestJson{}): {
				VendorId:  "iso15118",
				MessageId: "CertificateSigned",
			},
		},
	}

	return &Router{
		CallRoutes: map[string]handlers.CallRoute{
			"BootNotification": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.BootNotificationJson) },
				RequestSchema:  "ocpp16/BootNotification.json",
				ResponseSchema: "ocpp16/BootNotificationResponse.json",
				Handler: handlers16.BootNotificationHandler{
					Clock:               clk,
					RuntimeDetailsStore: engine,
					HeartbeatInterval:   int(heartbeatInterval.Seconds()),
				},
			},
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.HeartbeatJson) },
				RequestSchema:  "ocpp16/Heartbeat.json",
				ResponseSchema: "ocpp16/HeartbeatResponse.json",
				Handler: handlers16.HeartbeatHandler{
					Clock: clk,
				},
			},
			"StatusNotification": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.StatusNotificationJson) },
				RequestSchema:  "ocpp16/StatusNotification.json",
				ResponseSchema: "ocpp16/StatusNotificationResponse.json",
				Handler:        handlers.CallHandlerFunc(handlers16.StatusNotificationHandler),
			},
			"Authorize": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.AuthorizeJson) },
				RequestSchema:  "ocpp16/Authorize.json",
				ResponseSchema: "ocpp16/AuthorizeResponse.json",
				Handler: handlers16.AuthorizeHandler{
					TokenStore: engine,
				},
			},
			"StartTransaction": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.StartTransactionJson) },
				RequestSchema:  "ocpp16/StartTransaction.json",
				ResponseSchema: "ocpp16/StartTransactionResponse.json",
				Handler: handlers16.StartTransactionHandler{
					Clock:            clk,
					TokenStore:       engine,
					TransactionStore: engine,
				},
			},
			"StopTransaction": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.StopTransactionJson) },
				RequestSchema:  "ocpp16/StopTransaction.json",
				ResponseSchema: "ocpp16/StopTransactionResponse.json",
				Handler: handlers16.StopTransactionHandler{
					Clock:            clk,
					TokenStore:       engine,
					TransactionStore: engine,
				},
			},
			"MeterValues": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.MeterValuesJson) },
				RequestSchema:  "ocpp16/MeterValues.json",
				ResponseSchema: "ocpp16/MeterValuesResponse.json",
				Handler: handlers16.MeterValuesHandler{
					TransactionStore: engine,
				},
			},
			"DataTransfer": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.DataTransferJson) },
				RequestSchema:  "ocpp16/DataTransfer.json",
				ResponseSchema: "ocpp16/DataTransferResponse.json",
				Handler: handlers16.DataTransferHandler{
					SchemaFS: schemaFS,
					CallRoutes: map[string]map[string]handlers.CallRoute{
						"org.openchargealliance.iso15118pnc": {
							"Authorize": {
								NewRequest:     func() ocpp.Request { return new(ocpp201.AuthorizeRequestJson) },
								RequestSchema:  "ocpp201/AuthorizeRequest.json",
								ResponseSchema: "ocpp201/AuthorizeResponse.json",
								Handler: handlers201.AuthorizeHandler{
									TokenStore:                   engine,
									CertificateValidationService: certValidationService,
								},
							},
							"GetCertificateStatus": {
								NewRequest:     func() ocpp.Request { return new(ocpp201.GetCertificateStatusRequestJson) },
								RequestSchema:  "ocpp201/GetCertificateStatusRequest.json",
								ResponseSchema: "ocpp201/GetCertificateStatusResponse.json",
								Handler: handlers201.GetCertificateStatusHandler{
									CertificateValidationService: certValidationService,
								},
							},
							"SignCertificate": {
								NewRequest:     func() ocpp.Request { return new(ocpp201.SignCertificateRequestJson) },
								RequestSchema:  "ocpp201/SignCertificateRequest.json",
								ResponseSchema: "ocpp201/SignCertificateResponse.json",
								Handler: handlers201.SignCertificateHandler{
									ChargeStationCertificateProvider: chargeStationCertProvider,
									CallMaker:                        dataTransferCallMaker,
								},
							},
							"Get15118EVCertificate": {
								NewRequest:     func() ocpp.Request { return new(ocpp201.Get15118EVCertificateRequestJson) },
								RequestSchema:  "ocpp201/Get15118EVCertificateRequest.json",
								ResponseSchema: "ocpp201/Get15118EVCertificateResponse.json",
								Handler: handlers201.Get15118EvCertificateHandler{
									ContractCertificateProvider: contractCertProvider,
								},
							},
						},
						"iso15118": { // has2be extensions
							"Authorize": {
								NewRequest:     func() ocpp.Request { return new(has2be.AuthorizeRequestJson) },
								RequestSchema:  "has2be/AuthorizeRequest.json",
								ResponseSchema: "has2be/AuthorizeResponse.json",
								Handler: handlersHasToBe.AuthorizeHandler{
									Handler201: handlers201.AuthorizeHandler{
										TokenStore:                   engine,
										CertificateValidationService: certValidationService,
									},
								},
							},
							"GetCertificateStatus": {
								NewRequest:     func() ocpp.Request { return new(has2be.GetCertificateStatusRequestJson) },
								RequestSchema:  "has2be/GetCertificateStatusRequest.json",
								ResponseSchema: "has2be/GetCertificateStatusResponse.json",
								Handler: handlersHasToBe.GetCertificateStatusHandler{
									Handler201: handlers201.GetCertificateStatusHandler{
										CertificateValidationService: certValidationService,
									},
								},
							},
							"Get15118EVCertificate": {
								NewRequest:     func() ocpp.Request { return new(has2be.Get15118EVCertificateRequestJson) },
								RequestSchema:  "has2be/Get15118EVCertificateRequest.json",
								ResponseSchema: "has2be/Get15118EVCertificateResponse.json",
								Handler: handlersHasToBe.Get15118EvCertificateHandler{
									Handler201: handlers201.Get15118EvCertificateHandler{
										ContractCertificateProvider: contractCertProvider,
									},
								},
							},
							"SignCertificate": {
								NewRequest:     func() ocpp.Request { return new(has2be.SignCertificateRequestJson) },
								RequestSchema:  "has2be/SignCertificateRequestJson.json",
								ResponseSchema: "has2be/SignCertificateRequestJson.json",
								Handler: handlersHasToBe.SignCertificateHandler{
									Handler201: handlers201.SignCertificateHandler{
										ChargeStationCertificateProvider: chargeStationCertProvider,
										CallMaker:                        dataTransferCallMaker,
									},
								},
							},
						},
					},
				},
			},
		},
		CallResultRoutes: map[string]handlers.CallResultRoute{
			"DataTransfer": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.DataTransferJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp16.DataTransferResponseJson) },
				RequestSchema:  "ocpp16/DataTransfer.json",
				ResponseSchema: "ocpp16/DataTransferResponse.json",
				Handler: handlers16.DataTransferResultHandler{
					SchemaFS: schemaFS,
					CallResultRoutes: map[string]map[string]handlers.CallResultRoute{
						"org.openchargealliance.iso15118pnc": {
							"CertificateSigned": {
								NewRequest:     func() ocpp.Request { return new(ocpp201.CertificateSignedRequestJson) },
								NewResponse:    func() ocpp.Response { return new(ocpp201.CertificateSignedResponseJson) },
								RequestSchema:  "ocpp201/CertificateSignedRequest.json",
								ResponseSchema: "ocpp201/CertificateSignedResponse.json",
								Handler:        handlers201.CertificateSignedResultHandler{},
							},
						},
						"iso15118": { // has2be extensions
							"CertificateSigned": {
								NewRequest:     func() ocpp.Request { return new(has2be.CertificateSignedRequestJson) },
								NewResponse:    func() ocpp.Response { return new(has2be.CertificateSignedResponseJson) },
								RequestSchema:  "has2be/CertificateSignedRequest.json",
								ResponseSchema: "has2be/CertificateSignedResponse.json",
								Handler:        handlersHasToBe.CertificateSignedResultHandler{},
							},
						},
					},
				},
			},
			"ChangeConfiguration": {
				NewRequest:     func() ocpp.Request { return new(ocpp16.ChangeConfigurationJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp16.ChangeConfigurationResponseJson) },
				RequestSchema:  "ocpp16/ChangeConfiguration.json",
				ResponseSchema: "ocpp16/ChangeConfigurationResponse.json",
				Handler: handlers16.ChangeConfigurationResultHandler{
					SettingsStore: engine,
				},
			},
		},
	}
}

func NewV201Router(emitter Emitter,
	clk clock.PassiveClock,
	engine store.Engine,
	tariffService services.TariffService,
	certValidationService services.CertificateValidationService,
	chargeStationCertProvider services.ChargeStationCertificateProvider,
	contractCertProvider services.ContractCertificateProvider,
	heartbeatInterval time.Duration) *Router {

	callMaker := BasicCallMaker{
		E: emitter,
		Actions: map[reflect.Type]string{
			reflect.TypeOf(&ocpp201.CertificateSignedRequestJson{}): "CertificateSigned",
		},
	}

	return &Router{
		CallRoutes: map[string]handlers.CallRoute{
			"BootNotification": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.BootNotificationRequestJson) },
				RequestSchema:  "ocpp201/BootNotificationRequest.json",
				ResponseSchema: "ocpp201/BootNotificationResponse.json",
				Handler: handlers201.BootNotificationHandler{
					Clock:             clk,
					HeartbeatInterval: int(heartbeatInterval.Seconds()),
				},
			},
			"Heartbeat": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.HeartbeatRequestJson) },
				RequestSchema:  "ocpp201/HeartbeatRequest.json",
				ResponseSchema: "ocpp201/HeartbeatResponse.json",
				Handler: handlers201.HeartbeatHandler{
					Clock: clk,
				},
			},
			"StatusNotification": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.StatusNotificationRequestJson) },
				RequestSchema:  "ocpp201/StatusNotificationRequest.json",
				ResponseSchema: "ocpp201/StatusNotificationResponse.json",
				Handler:        handlers.CallHandlerFunc(handlers201.StatusNotificationHandler),
			},
			"Authorize": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.AuthorizeRequestJson) },
				RequestSchema:  "ocpp201/AuthorizeRequest.json",
				ResponseSchema: "ocpp201/AuthorizeResponse.json",
				Handler: handlers201.AuthorizeHandler{
					TokenStore:                   engine,
					CertificateValidationService: certValidationService,
				},
			},
			"TransactionEvent": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.TransactionEventRequestJson) },
				RequestSchema:  "ocpp201/TransactionEventRequest.json",
				ResponseSchema: "ocpp201/TransactionEventResponse.json",
				Handler: handlers201.TransactionEventHandler{
					TransactionStore: engine,
					TariffService:    tariffService,
				},
			},
			"GetCertificateStatus": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.GetCertificateStatusRequestJson) },
				RequestSchema:  "ocpp201/GetCertificateStatusRequest.json",
				ResponseSchema: "ocpp201/GetCertificateStatusResponse.json",
				Handler: handlers201.GetCertificateStatusHandler{
					CertificateValidationService: certValidationService,
				},
			},
			"SignCertificate": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.SignCertificateRequestJson) },
				RequestSchema:  "ocpp201/SignCertificateRequest.json",
				ResponseSchema: "ocpp201/SignCertificateResponse.json",
				Handler: handlers201.SignCertificateHandler{
					ChargeStationCertificateProvider: chargeStationCertProvider,
					CallMaker:                        callMaker,
				},
			},
			"Get15118EVCertificate": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.Get15118EVCertificateRequestJson) },
				RequestSchema:  "ocpp201/Get15118EVCertificateRequest.json",
				ResponseSchema: "ocpp201/Get15118EVCertificateResponse.json",
				Handler: handlers201.Get15118EvCertificateHandler{
					ContractCertificateProvider: contractCertProvider,
				},
			},
		},
		CallResultRoutes: map[string]handlers.CallResultRoute{
			"CertificateSigned": {
				NewRequest:     func() ocpp.Request { return new(ocpp201.CertificateSignedRequestJson) },
				NewResponse:    func() ocpp.Response { return new(ocpp201.CertificateSignedResponseJson) },
				RequestSchema:  "ocpp201/CertificateSignedRequest.json",
				ResponseSchema: "ocpp201/CertificateSignedResponse.json",
				Handler:        handlers201.CertificateSignedResultHandler{},
			},
		},
	}
}

func (r Router) Route(ctx context.Context, chargeStationId string, message Message, emitter Emitter, schemaFS fs.FS) error {
	switch message.MessageType {
	case MessageTypeCall:
		route, ok := r.CallRoutes[message.Action]
		if !ok {
			return fmt.Errorf("routing request: %w", NewError(ErrorNotImplemented, fmt.Errorf("%s not implemented", message.Action)))
		}
		err := schemas.Validate(message.RequestPayload, schemaFS, route.RequestSchema)
		if err != nil {
			var validationErr *jsonschema.ValidationError
			if errors.As(validationErr, &validationErr) {
				err = NewError(ErrorFormatViolation, err)
			}
			return fmt.Errorf("validating %s request: %w", message.Action, err)
		}
		req := route.NewRequest()
		err = json.Unmarshal(message.RequestPayload, &req)
		if err != nil {
			return fmt.Errorf("unmarshalling %s request payload: %w", message.Action, err)
		}
		resp, err := route.Handler.HandleCall(ctx, chargeStationId, req)
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("no response or error for %s", message.Action)
		}
		responseJson, err := json.Marshal(resp)
		if err != nil {
			return fmt.Errorf("marshalling %s call response: %w", message.Action, err)
		}
		err = schemas.Validate(responseJson, schemaFS, route.ResponseSchema)
		if err != nil {
			mqttErr := NewError(ErrorPropertyConstraintViolation, err)
			slog.Warn("response not valid", slog.String("action", message.Action), mqttErr)
		}
		out := &Message{
			MessageType:     MessageTypeCallResult,
			Action:          message.Action,
			MessageId:       message.MessageId,
			ResponsePayload: responseJson,
		}
		err = emitter.Emit(ctx, chargeStationId, out)
		if err != nil {
			return fmt.Errorf("sending call response: %w", err)
		}
	case MessageTypeCallResult:
		route, ok := r.CallResultRoutes[message.Action]
		if !ok {
			return fmt.Errorf("routing request: %w", NewError(ErrorNotImplemented, fmt.Errorf("%s result not implemented", message.Action)))
		}
		err := schemas.Validate(message.RequestPayload, schemaFS, route.RequestSchema)
		if err != nil {
			return fmt.Errorf("validating %s request: %w", message.Action, err)
		}
		err = schemas.Validate(message.ResponsePayload, schemaFS, route.ResponseSchema)
		if err != nil {
			var validationErr *jsonschema.ValidationError
			if errors.As(validationErr, &validationErr) {
				err = NewError(ErrorFormatViolation, err)
			}
			return fmt.Errorf("validating %s response: %w", message.Action, err)
		}
		req := route.NewRequest()
		err = json.Unmarshal(message.RequestPayload, &req)
		if err != nil {
			return fmt.Errorf("unmarshalling %s request payload: %w", message.Action, err)
		}
		resp := route.NewResponse()
		err = json.Unmarshal(message.ResponsePayload, &resp)
		if err != nil {
			return fmt.Errorf("unmarshalling %s response payload: %v", message.Action, err)
		}
		err = route.Handler.HandleCallResult(ctx, chargeStationId, req, resp, message.State)
		if err != nil {
			return err
		}
	case MessageTypeCallError:
		// TODO: what do we want to do with errors?
		return errors.New("we shouldn't get here at the moment")
	}

	return nil
}
