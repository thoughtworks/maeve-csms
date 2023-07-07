// SPDX-License-Identifier: Apache-2.0

package ocpp201

type Iso15118EVCertificateStatusEnumType string

const Iso15118EVCertificateStatusEnumTypeAccepted Iso15118EVCertificateStatusEnumType = "Accepted"
const Iso15118EVCertificateStatusEnumTypeFailed Iso15118EVCertificateStatusEnumType = "Failed"

type Get15118EVCertificateResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Raw CertificateInstallationRes response for the EV, Base64 encoded.
	//
	ExiResponse string `json:"exiResponse" yaml:"exiResponse" mapstructure:"exiResponse"`

	// Status corresponds to the JSON schema field "status".
	Status Iso15118EVCertificateStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*Get15118EVCertificateResponseJson) IsResponse() {}
