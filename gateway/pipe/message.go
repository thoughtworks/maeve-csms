// SPDX-License-Identifier: Apache-2.0

package pipe

import (
	"context"
	"encoding/json"
	"github.com/thoughtworks/maeve-csms/gateway/ocpp"
)

type GatewayMessage struct {
	Context          context.Context  `json:"-"`
	MessageType      ocpp.MessageType `json:"type"`
	Action           string           `json:"action"`
	MessageId        string           `json:"id"`
	RequestPayload   json.RawMessage  `json:"request,omitempty"`
	ResponsePayload  json.RawMessage  `json:"response,omitempty"`
	ErrorCode        ocpp.ErrorCode   `json:"error_code,omitempty"`
	ErrorDescription string           `json:"error_description,omitempty"`
	State            json.RawMessage  `json:"state,omitempty"`
}
