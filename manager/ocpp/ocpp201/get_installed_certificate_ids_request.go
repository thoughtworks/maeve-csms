// SPDX-License-Identifier: Apache-2.0

package ocpp201

type GetInstalledCertificateIdsRequestJson struct {
	// Indicates the type of certificates requested. When omitted, all certificate
	// types are requested.
	//
	CertificateType []GetCertificateIdUseEnumType `json:"certificateType,omitempty" yaml:"certificateType,omitempty" mapstructure:"certificateType,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

func (*GetInstalledCertificateIdsRequestJson) IsRequest() {}
