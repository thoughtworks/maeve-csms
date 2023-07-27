// SPDX-License-Identifier: Apache-2.0

package mqtt

import (
	"encoding/json"
)

type MessageType int

func (m MessageType) String() string {
	switch m {
	case MessageTypeCall:
		return "call"
	case MessageTypeCallResult:
		return "call_result"
	case MessageTypeCallError:
		return "call_error"
	default:
		return "unknown"
	}
}

const (
	MessageTypeCall MessageType = iota + 2
	MessageTypeCallResult
	MessageTypeCallError
)

type Message struct {
	MessageType      MessageType     `json:"type"`
	Action           string          `json:"action"`
	MessageId        string          `json:"id"`
	RequestPayload   json.RawMessage `json:"request,omitempty"`
	ResponsePayload  json.RawMessage `json:"response,omitempty"`
	ErrorCode        ErrorCode       `json:"error_code,omitempty"`
	ErrorDescription string          `json:"error_description,omitempty"`
	State            json.RawMessage `json:"state,omitempty"`
}

func NewErrorMessage(action, messageId string, code ErrorCode, err error) *Message {
	return &Message{
		MessageType:      MessageTypeCallError,
		Action:           action,
		MessageId:        messageId,
		ErrorCode:        code,
		ErrorDescription: err.Error(),
	}
}
