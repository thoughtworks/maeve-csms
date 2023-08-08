package has2be

type Get15118EVCertificateRequestJson struct {
	// Schema version currently used for the 15118 session between EV and Charging
	// Station. Needed for parsing of the EXI stream by the CSMS.
	A15118SchemaVersion *string `json:"15118SchemaVersion,omitempty" yaml:"15118SchemaVersion,omitempty" mapstructure:"15118SchemaVersion,omitempty"`

	// Raw CertificateInstallationReq request from EV, Base64 encoded.
	ExiRequest string `json:"exiRequest" yaml:"exiRequest" mapstructure:"exiRequest"`
}

func (*Get15118EVCertificateRequestJson) IsRequest() {}
