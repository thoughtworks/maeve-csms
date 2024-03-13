// SPDX-License-Identifier: Apache-2.0

package ocpp201

type ChangeAvailabilityStatusEnumType string

const ChangeAvailabilityStatusEnumTypeAccepted ChangeAvailabilityStatusEnumType = "Accepted"
const ChangeAvailabilityStatusEnumTypeRejected ChangeAvailabilityStatusEnumType = "Rejected"
const ChangeAvailabilityStatusEnumTypeScheduled ChangeAvailabilityStatusEnumType = "Scheduled"

type ChangeAvailabilityResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status ChangeAvailabilityStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*ChangeAvailabilityResponseJson) IsResponse() {}
