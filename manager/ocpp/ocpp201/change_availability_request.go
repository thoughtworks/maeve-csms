// SPDX-License-Identifier: Apache-2.0

package ocpp201

type OperationalStatusEnumType string

const OperationalStatusEnumTypeInoperative OperationalStatusEnumType = "Inoperative"
const OperationalStatusEnumTypeOperative OperationalStatusEnumType = "Operative"

type ChangeAvailabilityRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Evse corresponds to the JSON schema field "evse".
	Evse *EVSEType `json:"evse,omitempty" yaml:"evse,omitempty" mapstructure:"evse,omitempty"`

	// OperationalStatus corresponds to the JSON schema field "operationalStatus".
	OperationalStatus OperationalStatusEnumType `json:"operationalStatus" yaml:"operationalStatus" mapstructure:"operationalStatus"`
}

func (*ChangeAvailabilityRequestJson) IsRequest() {}
