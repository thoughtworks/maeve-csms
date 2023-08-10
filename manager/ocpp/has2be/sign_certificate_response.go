package has2be

type GenericStatusEnumType string

const GenericStatusEnumTypeAccepted GenericStatusEnumType = "Accepted"
const GenericStatusEnumTypeRejected GenericStatusEnumType = "Rejected"

type SignCertificateResponseJson struct {
	// Status corresponds to the JSON schema field "status".
	Status GenericStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`
}

func (*SignCertificateResponseJson) IsResponse() {}
