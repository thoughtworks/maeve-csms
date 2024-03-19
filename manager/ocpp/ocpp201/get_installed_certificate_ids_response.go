// SPDX-License-Identifier: Apache-2.0

package ocpp201

type CertificateHashDataChainType struct {
	// CertificateHashData corresponds to the JSON schema field "certificateHashData".
	CertificateHashData CertificateHashDataType `json:"certificateHashData" yaml:"certificateHashData" mapstructure:"certificateHashData"`

	// CertificateType corresponds to the JSON schema field "certificateType".
	CertificateType GetCertificateIdUseEnumType `json:"certificateType" yaml:"certificateType" mapstructure:"certificateType"`

	// ChildCertificateHashData corresponds to the JSON schema field
	// "childCertificateHashData".
	ChildCertificateHashData []CertificateHashDataType `json:"childCertificateHashData,omitempty" yaml:"childCertificateHashData,omitempty" mapstructure:"childCertificateHashData,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

type CertificateHashDataType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// HashAlgorithm corresponds to the JSON schema field "hashAlgorithm".
	HashAlgorithm HashAlgorithmEnumType `json:"hashAlgorithm" yaml:"hashAlgorithm" mapstructure:"hashAlgorithm"`

	// Hashed value of the issuers public key
	//
	IssuerKeyHash string `json:"issuerKeyHash" yaml:"issuerKeyHash" mapstructure:"issuerKeyHash"`

	// Hashed value of the Issuer DN (Distinguished Name).
	//
	//
	IssuerNameHash string `json:"issuerNameHash" yaml:"issuerNameHash" mapstructure:"issuerNameHash"`

	// The serial number of the certificate.
	//
	SerialNumber string `json:"serialNumber" yaml:"serialNumber" mapstructure:"serialNumber"`
}

type GetCertificateIdUseEnumType string

const GetCertificateIdUseEnumTypeCSMSRootCertificate GetCertificateIdUseEnumType = "CSMSRootCertificate"
const GetCertificateIdUseEnumTypeMORootCertificate GetCertificateIdUseEnumType = "MORootCertificate"
const GetCertificateIdUseEnumTypeManufacturerRootCertificate GetCertificateIdUseEnumType = "ManufacturerRootCertificate"
const GetCertificateIdUseEnumTypeV2GCertificateChain GetCertificateIdUseEnumType = "V2GCertificateChain"
const GetCertificateIdUseEnumTypeV2GRootCertificate GetCertificateIdUseEnumType = "V2GRootCertificate"

type GetInstalledCertificateIdsResponseJson struct {
	// CertificateHashDataChain corresponds to the JSON schema field
	// "certificateHashDataChain".
	CertificateHashDataChain []CertificateHashDataChainType `json:"certificateHashDataChain,omitempty" yaml:"certificateHashDataChain,omitempty" mapstructure:"certificateHashDataChain,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status GetInstalledCertificateStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*GetInstalledCertificateIdsResponseJson) IsResponse() {}

type GetInstalledCertificateStatusEnumType string

const GetInstalledCertificateStatusEnumTypeAccepted GetInstalledCertificateStatusEnumType = "Accepted"
const GetInstalledCertificateStatusEnumTypeNotFound GetInstalledCertificateStatusEnumType = "NotFound"
