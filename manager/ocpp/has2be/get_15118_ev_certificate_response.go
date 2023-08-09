package has2be

type Iso15118EVCertificateStatusEnumType string

const Iso15118EVCertificateStatusEnumTypeAccepted Iso15118EVCertificateStatusEnumType = "Accepted"
const Iso15118EVCertificateStatusEnumTypeFailed Iso15118EVCertificateStatusEnumType = "Failed"

type CertificateChainType struct {
	// Certificate corresponds to the JSON schema field "certificate".
	Certificate string `json:"certificate" yaml:"certificate" mapstructure:"certificate"`

	// ChildCertificate corresponds to the JSON schema field "childCertificate".
	ChildCertificate []string `json:"childCertificate,omitempty" yaml:"childCertificate,omitempty" mapstructure:"childCertificate,omitempty"`
}

type Get15118EVCertificateResponseJson struct {
	// ContractSignatureCertificateChain corresponds to the JSON schema field
	// "contractSignatureCertificateChain".
	ContractSignatureCertificateChain *CertificateChainType `json:"contractSignatureCertificateChain,omitempty" yaml:"contractSignatureCertificateChain,omitempty" mapstructure:"contractSignatureCertificateChain,omitempty"`

	// Raw CertificateInstallationRes response for the EV, Base64 encoded.
	ExiResponse string `json:"exiResponse" yaml:"exiResponse" mapstructure:"exiResponse"`

	// SaProvisioningCertificateChain corresponds to the JSON schema field
	// "saProvisioningCertificateChain".
	SaProvisioningCertificateChain *CertificateChainType `json:"saProvisioningCertificateChain,omitempty" yaml:"saProvisioningCertificateChain,omitempty" mapstructure:"saProvisioningCertificateChain,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status Iso15118EVCertificateStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`
}

func (*Get15118EVCertificateResponseJson) IsResponse() {}
