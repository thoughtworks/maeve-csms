package has2be

type IdTokenEnumType string

const IdTokenEnumTypeEMAID IdTokenEnumType = "eMAID"

// Contains a case insensitive identifier to use for the authorization and the type
// of authorization to support multiple forms of identifiers.
type IdTokenType struct {
	// IdToken is case insensitive. Will hold the eMAID.
	//
	IdToken string `json:"idToken" yaml:"idToken" mapstructure:"idToken"`

	// Type corresponds to the JSON schema field "type".
	Type IdTokenEnumType `json:"type" yaml:"type" mapstructure:"type"`
}

type AuthorizeRequestJson struct {
	// ISO15118CertificateHashData corresponds to the JSON schema field
	// "15118CertificateHashData".
	ISO15118CertificateHashData []OCSPRequestDataType `json:"15118CertificateHashData" yaml:"15118CertificateHashData" mapstructure:"15118CertificateHashData"`

	// IdToken corresponds to the JSON schema field "idToken".
	IdToken IdTokenType `json:"idToken" yaml:"idToken" mapstructure:"idToken"`
}

func (*AuthorizeRequestJson) IsRequest() {}
