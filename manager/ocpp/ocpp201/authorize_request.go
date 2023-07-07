// SPDX-License-Identifier: Apache-2.0

package ocpp201

type HashAlgorithmEnumType string

const HashAlgorithmEnumTypeSHA256 HashAlgorithmEnumType = "SHA256"
const HashAlgorithmEnumTypeSHA384 HashAlgorithmEnumType = "SHA384"
const HashAlgorithmEnumTypeSHA512 HashAlgorithmEnumType = "SHA512"

type OCSPRequestDataType struct {
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

	// This contains the responder URL (Case insensitive).
	//
	//
	ResponderURL string `json:"responderURL" yaml:"responderURL" mapstructure:"responderURL"`

	// The serial number of the certificate.
	//
	SerialNumber string `json:"serialNumber" yaml:"serialNumber" mapstructure:"serialNumber"`
}

type AuthorizeRequestJson struct {
	// The X.509 certificated presented by EV and encoded in PEM format.
	//
	Certificate *string `json:"certificate,omitempty" yaml:"certificate,omitempty" mapstructure:"certificate,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// IdToken corresponds to the JSON schema field "idToken".
	IdToken IdTokenType `json:"idToken" yaml:"idToken" mapstructure:"idToken"`

	// Iso15118CertificateHashData corresponds to the JSON schema field
	// "iso15118CertificateHashData".
	Iso15118CertificateHashData *[]OCSPRequestDataType `json:"iso15118CertificateHashData,omitempty" yaml:"iso15118CertificateHashData,omitempty" mapstructure:"iso15118CertificateHashData,omitempty"`
}

func (*AuthorizeRequestJson) IsRequest() {}
