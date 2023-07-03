package ocpp201

type CertificateSignedRequestJson struct {
	// The signed PEM encoded X.509 certificate. This can also contain the necessary
	// sub CA certificates. In that case, the order of the bundle should follow the
	// certificate chain, starting from the leaf certificate.
	//
	// The Configuration Variable
	// &lt;&lt;configkey-max-certificate-chain-size,MaxCertificateChainSize&gt;&gt;
	// can be used to limit the maximum size of this field.
	//
	CertificateChain string `json:"certificateChain" yaml:"certificateChain" mapstructure:"certificateChain"`

	// CertificateType corresponds to the JSON schema field "certificateType".
	CertificateType *CertificateSigningUseEnumType `json:"certificateType,omitempty" yaml:"certificateType,omitempty" mapstructure:"certificateType,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

func (*CertificateSignedRequestJson) IsRequest() {}
