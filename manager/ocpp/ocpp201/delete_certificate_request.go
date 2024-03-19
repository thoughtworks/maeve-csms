// SPDX-License-Identifier: Apache-2.0

package ocpp201

type DeleteCertificateRequestJson struct {
	// CertificateHashData corresponds to the JSON schema field "certificateHashData".
	CertificateHashData CertificateHashDataType `json:"certificateHashData" yaml:"certificateHashData" mapstructure:"certificateHashData"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

func (*DeleteCertificateRequestJson) IsRequest() {}
