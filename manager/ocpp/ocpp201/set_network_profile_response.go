// SPDX-License-Identifier: Apache-2.0

package ocpp201

type SetNetworkProfileStatusEnumType string

const SetNetworkProfileStatusEnumTypeAccepted SetNetworkProfileStatusEnumType = "Accepted"
const SetNetworkProfileStatusEnumTypeFailed SetNetworkProfileStatusEnumType = "Failed"
const SetNetworkProfileStatusEnumTypeRejected SetNetworkProfileStatusEnumType = "Rejected"

type SetNetworkProfileResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status SetNetworkProfileStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*SetNetworkProfileResponseJson) IsResponse() {}
