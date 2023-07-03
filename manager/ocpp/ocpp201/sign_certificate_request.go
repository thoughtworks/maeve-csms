package ocpp201

type CertificateSigningUseEnumType string

const CertificateSigningUseEnumTypeChargingStationCertificate CertificateSigningUseEnumType = "ChargingStationCertificate"
const CertificateSigningUseEnumTypeV2GCertificate CertificateSigningUseEnumType = "V2GCertificate"

type SignCertificateRequestJson struct {
	// CertificateType corresponds to the JSON schema field "certificateType".
	CertificateType *CertificateSigningUseEnumType `json:"certificateType,omitempty" yaml:"certificateType,omitempty" mapstructure:"certificateType,omitempty"`

	// The Charging Station SHALL send the public key in form of a Certificate Signing
	// Request (CSR) as described in RFC 2986 [22] and then PEM encoded, using the
	// &lt;&lt;signcertificaterequest,SignCertificateRequest&gt;&gt; message.
	//
	Csr string `json:"csr" yaml:"csr" mapstructure:"csr"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

func (*SignCertificateRequestJson) IsRequest() {}
