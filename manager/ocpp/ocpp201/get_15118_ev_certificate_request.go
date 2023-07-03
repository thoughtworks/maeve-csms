package ocpp201

type CertificateActionEnumType string

const CertificateActionEnumTypeInstall CertificateActionEnumType = "Install"
const CertificateActionEnumTypeUpdate CertificateActionEnumType = "Update"

type Get15118EVCertificateRequestJson struct {
	// Action corresponds to the JSON schema field "action".
	Action CertificateActionEnumType `json:"action" yaml:"action" mapstructure:"action"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	ExiRequest string `json:"exiRequest" yaml:"exiRequest" mapstructure:"exiRequest"`

	Iso15118SchemaVersion string `json:"iso15118SchemaVersion" yaml:"iso15118SchemaVersion" mapstructure:"iso15118SchemaVersion"`
}

func (*Get15118EVCertificateRequestJson) IsRequest() {}
