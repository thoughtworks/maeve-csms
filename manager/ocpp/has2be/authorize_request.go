package has2be

type HashAlgorithmEnumType string

const HashAlgorithmEnumTypeSHA256 HashAlgorithmEnumType = "SHA256"
const HashAlgorithmEnumTypeSHA384 HashAlgorithmEnumType = "SHA384"
const HashAlgorithmEnumTypeSHA512 HashAlgorithmEnumType = "SHA512"

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

type OCSPRequestDataType struct {
	// HashAlgorithm corresponds to the JSON schema field "hashAlgorithm".
	HashAlgorithm HashAlgorithmEnumType `json:"hashAlgorithm" yaml:"hashAlgorithm" mapstructure:"hashAlgorithm"`

	// Hashed value of the issuers public key
	//
	IssuerKeyHash string `json:"issuerKeyHash" yaml:"issuerKeyHash" mapstructure:"issuerKeyHash"`

	// Hashed value of the Issuer DN (Distinguished Name).
	//
	//
	IssuerNameHash string `json:"issuerNameHash" yaml:"issuerNameHash" mapstructure:"issuerNameHash"`

	// This contains the responder URL (Case insensitive).
	//
	//
	ResponderURL *string `json:"responderURL,omitempty" yaml:"responderURL,omitempty" mapstructure:"responderURL,omitempty"`

	// The serial number of the certificate.
	//
	SerialNumber string `json:"serialNumber" yaml:"serialNumber" mapstructure:"serialNumber"`
}

type AuthorizeRequestJson struct {
	// ISO15118CertificateHashData corresponds to the JSON schema field
	// "15118CertificateHashData".
	ISO15118CertificateHashData []OCSPRequestDataType `json:"15118CertificateHashData" yaml:"15118CertificateHashData" mapstructure:"15118CertificateHashData"`

	// IdToken corresponds to the JSON schema field "idToken".
	IdToken IdTokenType `json:"idToken" yaml:"idToken" mapstructure:"idToken"`
}

func (*AuthorizeRequestJson) IsRequest() {}
