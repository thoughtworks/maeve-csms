package has2be

type CertificateSignedStatusEnumType string

type CertificateSignedResponseJson struct {
	// Status corresponds to the JSON schema field "status".
	Status CertificateSignedStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`
}

const CertificateSignedStatusEnumTypeAccepted CertificateSignedStatusEnumType = "Accepted"
const CertificateSignedStatusEnumTypeRejected CertificateSignedStatusEnumType = "Rejected"

func (*CertificateSignedResponseJson) IsResponse() {}
