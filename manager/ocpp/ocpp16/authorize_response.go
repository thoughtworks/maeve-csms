package ocpp16

type AuthorizeResponseJsonIdTagInfoStatus string

var enumValues_AuthorizeResponseJsonIdTagInfoStatus = []interface{}{
	"Accepted",
	"Blocked",
	"Expired",
	"Invalid",
	"ConcurrentTx",
}

type AuthorizeResponseJsonIdTagInfo struct {
	// ExpiryDate corresponds to the JSON schema field "expiryDate".
	ExpiryDate *string `json:"expiryDate,omitempty" yaml:"expiryDate,omitempty" mapstructure:"expiryDate,omitempty"`

	// ParentIdTag corresponds to the JSON schema field "parentIdTag".
	ParentIdTag *string `json:"parentIdTag,omitempty" yaml:"parentIdTag,omitempty" mapstructure:"parentIdTag,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status AuthorizeResponseJsonIdTagInfoStatus `json:"status" yaml:"status" mapstructure:"status"`
}

const AuthorizeResponseJsonIdTagInfoStatusAccepted AuthorizeResponseJsonIdTagInfoStatus = "Accepted"
const AuthorizeResponseJsonIdTagInfoStatusBlocked AuthorizeResponseJsonIdTagInfoStatus = "Blocked"
const AuthorizeResponseJsonIdTagInfoStatusConcurrentTx AuthorizeResponseJsonIdTagInfoStatus = "ConcurrentTx"
const AuthorizeResponseJsonIdTagInfoStatusExpired AuthorizeResponseJsonIdTagInfoStatus = "Expired"
const AuthorizeResponseJsonIdTagInfoStatusInvalid AuthorizeResponseJsonIdTagInfoStatus = "Invalid"

type AuthorizeResponseJson struct {
	// IdTagInfo corresponds to the JSON schema field "idTagInfo".
	IdTagInfo AuthorizeResponseJsonIdTagInfo `json:"idTagInfo" yaml:"idTagInfo" mapstructure:"idTagInfo"`
}

func (*AuthorizeResponseJson) IsResponse() {}
