// SPDX-License-Identifier: Apache-2.0

package ocpp201

type MessageTriggerEnumType string

const MessageTriggerEnumTypeBootNotification MessageTriggerEnumType = "BootNotification"
const MessageTriggerEnumTypeFirmwareStatusNotification MessageTriggerEnumType = "FirmwareStatusNotification"
const MessageTriggerEnumTypeHeartbeat MessageTriggerEnumType = "Heartbeat"
const MessageTriggerEnumTypeLogStatusNotification MessageTriggerEnumType = "LogStatusNotification"
const MessageTriggerEnumTypeMeterValues MessageTriggerEnumType = "MeterValues"
const MessageTriggerEnumTypePublishFirmwareStatusNotification MessageTriggerEnumType = "PublishFirmwareStatusNotification"
const MessageTriggerEnumTypeSignChargingStationCertificate MessageTriggerEnumType = "SignChargingStationCertificate"
const MessageTriggerEnumTypeSignCombinedCertificate MessageTriggerEnumType = "SignCombinedCertificate"
const MessageTriggerEnumTypeSignV2GCertificate MessageTriggerEnumType = "SignV2GCertificate"
const MessageTriggerEnumTypeStatusNotification MessageTriggerEnumType = "StatusNotification"
const MessageTriggerEnumTypeTransactionEvent MessageTriggerEnumType = "TransactionEvent"

type TriggerMessageRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Evse corresponds to the JSON schema field "evse".
	Evse *EVSEType `json:"evse,omitempty" yaml:"evse,omitempty" mapstructure:"evse,omitempty"`

	// RequestedMessage corresponds to the JSON schema field "requestedMessage".
	RequestedMessage MessageTriggerEnumType `json:"requestedMessage" yaml:"requestedMessage" mapstructure:"requestedMessage"`
}

func (*TriggerMessageRequestJson) IsRequest() {}
