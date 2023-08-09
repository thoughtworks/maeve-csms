package has2be

type GetCertificateStatusRequestJson struct {
	// OcspRequestData corresponds to the JSON schema field "ocspRequestData".
	OcspRequestData OCSPRequestDataType `json:"ocspRequestData" yaml:"ocspRequestData" mapstructure:"ocspRequestData"`
}

func (*GetCertificateStatusRequestJson) IsRequest() {}
