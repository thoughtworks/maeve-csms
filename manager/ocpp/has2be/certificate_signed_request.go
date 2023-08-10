package has2be

type CertificateSignedRequestJson struct {
	// The signed X.509 certificate, first DER encoded into binary, and then hex
	// encoded into a case insensitive string. This can also contain the necessary sub
	// CA certificates. In that case, the order should follow the certificate chain,
	// starting from the leaf certificate.

	// TypeOfCertificate corresponds to the JSON schema field "typeOfCertificate".
	TypeOfCertificate CertificateSigningUseEnumType `json:"typeOfCertificate" yaml:"typeOfCertificate" mapstructure:"typeOfCertificate"`
}

func (*CertificateSignedRequestJson) IsRequest() {}
