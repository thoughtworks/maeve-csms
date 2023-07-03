package ocpp201

type GetCertificateStatusRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// OcspRequestData corresponds to the JSON schema field "ocspRequestData".
	OcspRequestData OCSPRequestDataType `json:"ocspRequestData" yaml:"ocspRequestData" mapstructure:"ocspRequestData"`
}

func (*GetCertificateStatusRequestJson) IsRequest() {}
