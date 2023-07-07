// SPDX-License-Identifier: Apache-2.0

package ocpp201

type AuthorizeCertificateStatusEnumType string

const AuthorizeCertificateStatusEnumTypeAccepted AuthorizeCertificateStatusEnumType = "Accepted"
const AuthorizeCertificateStatusEnumTypeCertChainError AuthorizeCertificateStatusEnumType = "CertChainError"
const AuthorizeCertificateStatusEnumTypeCertificateExpired AuthorizeCertificateStatusEnumType = "CertificateExpired"
const AuthorizeCertificateStatusEnumTypeCertificateRevoked AuthorizeCertificateStatusEnumType = "CertificateRevoked"
const AuthorizeCertificateStatusEnumTypeContractCancelled AuthorizeCertificateStatusEnumType = "ContractCancelled"
const AuthorizeCertificateStatusEnumTypeNoCertificateAvailable AuthorizeCertificateStatusEnumType = "NoCertificateAvailable"
const AuthorizeCertificateStatusEnumTypeSignatureError AuthorizeCertificateStatusEnumType = "SignatureError"

type AuthorizeResponseJson struct {
	// CertificateStatus corresponds to the JSON schema field "certificateStatus".
	CertificateStatus *AuthorizeCertificateStatusEnumType `json:"certificateStatus,omitempty" yaml:"certificateStatus,omitempty" mapstructure:"certificateStatus,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// IdTokenInfo corresponds to the JSON schema field "idTokenInfo".
	IdTokenInfo IdTokenInfoType `json:"idTokenInfo" yaml:"idTokenInfo" mapstructure:"idTokenInfo"`
}

func (*AuthorizeResponseJson) IsResponse() {}
