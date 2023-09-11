// SPDX-License-Identifier: Apache-2.0

package ocpp16

type TriggerMessageJson struct {
	// ConnectorId corresponds to the JSON schema field "connectorId".
	ConnectorId *int `json:"connectorId,omitempty" yaml:"connectorId,omitempty" mapstructure:"connectorId,omitempty"`

	// RequestedMessage corresponds to the JSON schema field "requestedMessage".
	RequestedMessage TriggerMessageJsonRequestedMessage `json:"requestedMessage" yaml:"requestedMessage" mapstructure:"requestedMessage"`
}

func (*TriggerMessageJson) IsRequest() {}

type TriggerMessageJsonRequestedMessage string

const TriggerMessageJsonRequestedMessageBootNotification TriggerMessageJsonRequestedMessage = "BootNotification"
const TriggerMessageJsonRequestedMessageDiagnosticsStatusNotification TriggerMessageJsonRequestedMessage = "DiagnosticsStatusNotification"
const TriggerMessageJsonRequestedMessageFirmwareStatusNotification TriggerMessageJsonRequestedMessage = "FirmwareStatusNotification"
const TriggerMessageJsonRequestedMessageHeartbeat TriggerMessageJsonRequestedMessage = "Heartbeat"
const TriggerMessageJsonRequestedMessageMeterValues TriggerMessageJsonRequestedMessage = "MeterValues"
const TriggerMessageJsonRequestedMessageStatusNotification TriggerMessageJsonRequestedMessage = "StatusNotification"
