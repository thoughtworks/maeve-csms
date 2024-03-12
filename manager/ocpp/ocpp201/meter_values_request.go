// SPDX-License-Identifier: Apache-2.0

package ocpp201

// Request_ Body
// urn:x-enexis:ecdm:uid:2:234744
type MeterValuesRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Request_ Body. EVSEID. Numeric_ Identifier
	// urn:x-enexis:ecdm:uid:1:571101
	// This contains a number (&gt;0) designating an EVSE of the Charging Station. ‘0’
	// (zero) is used to designate the main power meter.
	//
	EvseId int `json:"evseId" yaml:"evseId" mapstructure:"evseId"`

	// MeterValue corresponds to the JSON schema field "meterValue".
	MeterValue []MeterValueType `json:"meterValue" yaml:"meterValue" mapstructure:"meterValue"`
}

func (*MeterValuesRequestJson) IsRequest() {}
