// SPDX-License-Identifier: Apache-2.0

package ocpp201

type ResetEnumType string

const ResetEnumTypeImmediate ResetEnumType = "Immediate"
const ResetEnumTypeOnIdle ResetEnumType = "OnIdle"

type ResetRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// This contains the ID of a specific EVSE that needs to be reset, instead of the
	// entire Charging Station.
	//
	EvseId *int `json:"evseId,omitempty" yaml:"evseId,omitempty" mapstructure:"evseId,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type ResetEnumType `json:"type" yaml:"type" mapstructure:"type"`
}

func (*ResetRequestJson) IsRequest() {}
