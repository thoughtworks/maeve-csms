package pipe

import (
	"encoding/json"
	"github.com/twlabs/ocpp2-broker-core/gateway/ocpp"
)

type GatewayMessage struct {
	MessageType      ocpp.MessageType `json:"type"`
	Action           string           `json:"action"`
	MessageId        string           `json:"id"`
	RequestPayload   json.RawMessage  `json:"request,omitempty"`
	ResponsePayload  json.RawMessage  `json:"response,omitempty"`
	ErrorCode        ocpp.ErrorCode   `json:"error_code,omitempty"`
	ErrorDescription string           `json:"error_description,omitempty"`
	State            json.RawMessage  `json:"state,omitempty"`
}
