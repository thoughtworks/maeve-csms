// SPDX-License-Identifier: Apache-2.0

package ocpp201

type UnlockConnectorResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status UnlockStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*UnlockConnectorResponseJson) IsResponse() {}

type UnlockStatusEnumType string

const UnlockStatusEnumTypeOngoingAuthorizedTransaction UnlockStatusEnumType = "OngoingAuthorizedTransaction"
const UnlockStatusEnumTypeUnknownConnector UnlockStatusEnumType = "UnknownConnector"
const UnlockStatusEnumTypeUnlockFailed UnlockStatusEnumType = "UnlockFailed"
const UnlockStatusEnumTypeUnlocked UnlockStatusEnumType = "Unlocked"
