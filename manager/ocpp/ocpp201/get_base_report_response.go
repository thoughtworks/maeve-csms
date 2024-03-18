// SPDX-License-Identifier: Apache-2.0

package ocpp201

type GenericDeviceModelStatusEnumType string

const GenericDeviceModelStatusEnumTypeAccepted GenericDeviceModelStatusEnumType = "Accepted"
const GenericDeviceModelStatusEnumTypeEmptyResultSet GenericDeviceModelStatusEnumType = "EmptyResultSet"
const GenericDeviceModelStatusEnumTypeNotSupported GenericDeviceModelStatusEnumType = "NotSupported"
const GenericDeviceModelStatusEnumTypeRejected GenericDeviceModelStatusEnumType = "Rejected"

type GetBaseReportResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status GenericDeviceModelStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*GetBaseReportResponseJson) IsResponse() {}
