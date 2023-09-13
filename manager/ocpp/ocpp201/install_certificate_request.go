// SPDX-License-Identifier: Apache-2.0

package ocpp201

type InstallCertificateUseEnumType string

const InstallCertificateUseEnumTypeCSMSRootCertificate InstallCertificateUseEnumType = "CSMSRootCertificate"
const InstallCertificateUseEnumTypeMORootCertificate InstallCertificateUseEnumType = "MORootCertificate"
const InstallCertificateUseEnumTypeManufacturerRootCertificate InstallCertificateUseEnumType = "ManufacturerRootCertificate"
const InstallCertificateUseEnumTypeV2GRootCertificate InstallCertificateUseEnumType = "V2GRootCertificate"

type InstallCertificateRequestJson struct {
	// A PEM encoded X.509 certificate.
	//
	Certificate string `json:"certificate" yaml:"certificate" mapstructure:"certificate"`

	// CertificateType corresponds to the JSON schema field "certificateType".
	CertificateType InstallCertificateUseEnumType `json:"certificateType" yaml:"certificateType" mapstructure:"certificateType"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

func (*InstallCertificateRequestJson) IsRequest() {}
