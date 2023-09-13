// SPDX-License-Identifier: Apache-2.0

package ocpp201

type InstallCertificateStatusEnumType string

const InstallCertificateStatusEnumTypeAccepted InstallCertificateStatusEnumType = "Accepted"
const InstallCertificateStatusEnumTypeFailed InstallCertificateStatusEnumType = "Failed"
const InstallCertificateStatusEnumTypeRejected InstallCertificateStatusEnumType = "Rejected"

type InstallCertificateResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status InstallCertificateStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*InstallCertificateResponseJson) IsResponse() {}
