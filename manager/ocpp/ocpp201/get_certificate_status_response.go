package ocpp201

type GetCertificateStatusEnumType string

const GetCertificateStatusEnumTypeAccepted GetCertificateStatusEnumType = "Accepted"
const GetCertificateStatusEnumTypeFailed GetCertificateStatusEnumType = "Failed"

type GetCertificateStatusResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// OCSPResponse class as defined in &lt;&lt;ref-ocpp_security_24, IETF RFC
	// 6960&gt;&gt;. DER encoded (as defined in &lt;&lt;ref-ocpp_security_24, IETF RFC
	// 6960&gt;&gt;), and then base64 encoded. MAY only be omitted when status is not
	// Accepted.
	//
	OcspResult *string `json:"ocspResult,omitempty" yaml:"ocspResult,omitempty" mapstructure:"ocspResult,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status GetCertificateStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*GetCertificateStatusResponseJson) IsResponse() {}
