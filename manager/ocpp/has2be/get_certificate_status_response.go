package has2be

type GetCertificateStatusEnumType string

const GetCertificateStatusEnumTypeAccepted GetCertificateStatusEnumType = "Accepted"
const GetCertificateStatusEnumTypeFailed GetCertificateStatusEnumType = "Failed"

type GetCertificateStatusResponseJson struct {
	// OCSPResponse class as defined in <<ref-ocpp_security_24, IETF RFC
	// 6960>>. DER encoded (as defined in <<ref-ocpp_security_24, IETF RFC
	// 6960>>;), and then base64 encoded. MAY only be omitted when status is not
	// Accepted.
	OcspResult *string `json:"ocspResult,omitempty" yaml:"ocspResult,omitempty" mapstructure:"ocspResult,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status GetCertificateStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`
}

func (*GetCertificateStatusResponseJson) IsResponse() {}
