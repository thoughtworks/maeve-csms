package has2be

type CertificateSigningUseEnumType string

const CertificateSigningUseEnumTypeChargingStationCertificate CertificateSigningUseEnumType = "ChargingStationCertificate"
const CertificateSigningUseEnumTypeV2GCertificate CertificateSigningUseEnumType = "V2GCertificate"

type SignCertificateRequestJson struct {
	// The Charging Station SHALL send the public key in form of a Certificate Signing
	// Request (CSR) as described in the X.509 standard.
	Csr string `json:"csr" yaml:"csr" mapstructure:"csr"`

	// TypeOfCertificate corresponds to the JSON schema field "typeOfCertificate".
	TypeOfCertificate *CertificateSigningUseEnumType `json:"typeOfCertificate,omitempty" yaml:"typeOfCertificate,omitempty" mapstructure:"typeOfCertificate,omitempty"`
}

func (*SignCertificateRequestJson) IsRequest() {}
