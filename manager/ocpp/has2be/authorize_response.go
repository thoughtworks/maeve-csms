package has2be

type AuthorizationStatusEnumType string

const AuthorizationStatusEnumTypeAccepted AuthorizationStatusEnumType = "Accepted"
const AuthorizationStatusEnumTypeBlocked AuthorizationStatusEnumType = "Blocked"
const AuthorizationStatusEnumTypeConcurrentTx AuthorizationStatusEnumType = "ConcurrentTx"
const AuthorizationStatusEnumTypeExpired AuthorizationStatusEnumType = "Expired"
const AuthorizationStatusEnumTypeInvalid AuthorizationStatusEnumType = "Invalid"
const AuthorizationStatusEnumTypeNoCredit AuthorizationStatusEnumType = "NoCredit"
const AuthorizationStatusEnumTypeNotAllowedTypeEVSE AuthorizationStatusEnumType = "NotAllowedTypeEVSE"
const AuthorizationStatusEnumTypeNotAtThisLocation AuthorizationStatusEnumType = "NotAtThisLocation"
const AuthorizationStatusEnumTypeNotAtThisTime AuthorizationStatusEnumType = "NotAtThisTime"
const AuthorizationStatusEnumTypeUnknown AuthorizationStatusEnumType = "Unknown"

type AuthorizeCertificateStatusEnumType string

const AuthorizeCertificateStatusEnumTypeAccepted AuthorizeCertificateStatusEnumType = "Accepted"
const AuthorizeCertificateStatusEnumTypeCertificateRevoked AuthorizeCertificateStatusEnumType = "CertificateRevoked"

// ID_ Token
// urn:x-oca:ocpp:uid:2:233247
// Contains status information about an identifier.
// It is advised to not stop charging for a token that expires during charging, as
// ExpiryDate is only used for caching purposes. If ExpiryDate is not given, the
// status has no end date.
type IdTokenInfoType struct {
	// ID_ Token. Expiry. Date_ Time
	// urn:x-oca:ocpp:uid:1:569373
	// Date and Time after which the token must be considered invalid.
	//
	CacheExpiryDateTime *string `json:"cacheExpiryDateTime,omitempty" yaml:"cacheExpiryDateTime,omitempty" mapstructure:"cacheExpiryDateTime,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status AuthorizationStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`
}

type AuthorizeResponseJson struct {
	// CertificateStatus corresponds to the JSON schema field "certificateStatus".
	CertificateStatus AuthorizeCertificateStatusEnumType `json:"certificateStatus" yaml:"certificateStatus" mapstructure:"certificateStatus"`

	// IdTokenInfo corresponds to the JSON schema field "idTokenInfo".
	IdTokenInfo IdTokenInfoType `json:"idTokenInfo" yaml:"idTokenInfo" mapstructure:"idTokenInfo"`
}

func (*AuthorizeResponseJson) IsResponse() {}
