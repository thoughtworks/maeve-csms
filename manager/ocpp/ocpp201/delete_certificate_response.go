// SPDX-License-Identifier: Apache-2.0

package ocpp201

type DeleteCertificateStatusEnumType string

const DeleteCertificateStatusEnumTypeAccepted DeleteCertificateStatusEnumType = "Accepted"
const DeleteCertificateStatusEnumTypeFailed DeleteCertificateStatusEnumType = "Failed"
const DeleteCertificateStatusEnumTypeNotFound DeleteCertificateStatusEnumType = "NotFound"

type DeleteCertificateResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status DeleteCertificateStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*DeleteCertificateResponseJson) IsResponse() {}
