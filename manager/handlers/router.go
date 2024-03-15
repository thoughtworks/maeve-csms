// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/thoughtworks/maeve-csms/manager/schemas"
	"github.com/thoughtworks/maeve-csms/manager/transport"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
	"io/fs"
)

// Router is the primary implementation of the transport.Router interface.
type Router struct {
	Emitter          transport.Emitter          // used to send responses to the gateway
	SchemaFS         fs.FS                      // used to obtain schema files
	OcppVersion      transport.OcppVersion      // the OCPP version that this router supports
	CallRoutes       map[string]CallRoute       // the set of routes for incoming calls (indexed by action)
	CallResultRoutes map[string]CallResultRoute // the set of routes for call results (indexed by action)
}

func (r Router) Handle(ctx context.Context, chargeStationId string, msg *transport.Message) {
	span := trace.SpanFromContext(ctx)

	err := r.route(ctx, chargeStationId, msg)
	if err != nil {
		slog.Error("unable to route message", slog.String("chargeStationId", chargeStationId), slog.String("action", msg.Action), "err", err)
		span.SetStatus(codes.Error, "routing request failed")
		span.RecordError(err)

		// only emit an error on a call (the charge station will not be expecting any response message)
		if msg.MessageType == transport.MessageTypeCall {
			var mqttError *transport.Error
			var errMsg *transport.Message
			if errors.As(err, &mqttError) {
				errMsg = transport.NewErrorMessage(msg.Action, msg.MessageId, mqttError.ErrorCode, mqttError.WrappedError)
			} else {
				errMsg = transport.NewErrorMessage(msg.Action, msg.MessageId, transport.ErrorInternalError, err)
			}
			err = r.Emitter.Emit(ctx, r.OcppVersion, chargeStationId, errMsg)
			if err != nil {
				slog.Error("unable to emit error message", "err", err)
			}
		}
	} else {
		span.SetStatus(codes.Ok, "ok")
	}
}

func (r Router) route(ctx context.Context, chargeStationId string, message *transport.Message) error {
	switch message.MessageType {
	case transport.MessageTypeCall:
		route, ok := r.CallRoutes[message.Action]
		if !ok {
			return fmt.Errorf("routing request: %w", transport.NewError(transport.ErrorNotImplemented, fmt.Errorf("%s not implemented", message.Action)))
		}
		err := schemas.Validate(message.RequestPayload, r.SchemaFS, route.RequestSchema)
		if err != nil {
			var validationErr *jsonschema.ValidationError
			if errors.As(validationErr, &validationErr) {
				err = transport.NewError(transport.ErrorFormatViolation, err)
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
		err = schemas.Validate(responseJson, r.SchemaFS, route.ResponseSchema)
		if err != nil {
			mqttErr := transport.NewError(transport.ErrorPropertyConstraintViolation, err)
			slog.Warn("response not valid", slog.String("action", message.Action), mqttErr)
		}
		out := &transport.Message{
			MessageType:     transport.MessageTypeCallResult,
			Action:          message.Action,
			MessageId:       message.MessageId,
			ResponsePayload: responseJson,
		}
		err = r.Emitter.Emit(ctx, r.OcppVersion, chargeStationId, out)
		if err != nil {
			return fmt.Errorf("sending call response: %w", err)
		}
	case transport.MessageTypeCallResult:
		route, ok := r.CallResultRoutes[message.Action]
		if !ok {
			return fmt.Errorf("routing request: %w", transport.NewError(transport.ErrorNotImplemented, fmt.Errorf("%s result not implemented", message.Action)))
		}
		err := schemas.Validate(message.RequestPayload, r.SchemaFS, route.RequestSchema)
		if err != nil {
			return fmt.Errorf("validating %s request: %w", message.Action, err)
		}
		err = schemas.Validate(message.ResponsePayload, r.SchemaFS, route.ResponseSchema)
		if err != nil {
			var validationErr *jsonschema.ValidationError
			if errors.As(validationErr, &validationErr) {
				err = transport.NewError(transport.ErrorFormatViolation, err)
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
	case transport.MessageTypeCallError:
		// TODO: what do we want to do with errors?
		return errors.New("we shouldn't get here at the moment")
	}

	return nil
}
