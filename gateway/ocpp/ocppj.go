// SPDX-License-Identifier: Apache-2.0

package ocpp

import (
	"encoding/json"
	"errors"
)

type MessageType int

const (
	MessageTypeCall MessageType = iota + 2
	MessageTypeCallResult
	MessageTypeCallError
)

type Message struct {
	MessageTypeId MessageType
	MessageId     string
	Data          []json.RawMessage
}

func (m Message) MarshalJSON() ([]byte, error) {
	var message []json.RawMessage

	data, err := json.Marshal(m.MessageTypeId)
	if err != nil {
		return nil, err
	}
	message = append(message, data)

	data, err = json.Marshal(m.MessageId)
	if err != nil {
		return nil, err
	}
	message = append(message, data)

	message = append(message, m.Data...)

	return json.Marshal(message)
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var message []json.RawMessage
	err := json.Unmarshal(data, &message)
	if err != nil {
		return err
	}

	ml := len(message)
	if ml < 1 {
		return errors.New("no message type")
	}

	err = json.Unmarshal(message[0], &m.MessageTypeId)
	if err != nil {
		return err
	}

	if ml > 1 {
		err = json.Unmarshal(message[1], &m.MessageId)
		if err != nil {
			return err
		}
	}

	if ml > 2 {
		m.Data = message[2:]
	}

	return nil
}

type ErrorCode string

const (
	ErrorFormatViolation               ErrorCode = "FormatViolation"
	ErrorGenericError                  ErrorCode = "GenericError"
	ErrorInternalError                 ErrorCode = "InternalError"
	ErrorMessageTypeNotSupported       ErrorCode = "MessageTypeNotSupported"
	ErrorNotImplemented                ErrorCode = "NotImplemented"
	ErrorNotSupported                  ErrorCode = "NotSupported"
	ErrorOccurrenceConstraintViolation ErrorCode = "OccurrenceConstraintViolation"
	ErrorPropertyConstraintViolation   ErrorCode = "PropertyConstraintViolation"
	ErrorProtocolError                 ErrorCode = "ProtocolError"
	ErrorRpcFrameworkError             ErrorCode = "RpcFrameworkError"
	ErrorSecurityError                 ErrorCode = "SecurityError"
	ErrorTypeConstraintViolation       ErrorCode = "TypeConstraintViolation"
)
