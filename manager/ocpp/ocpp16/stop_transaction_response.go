package ocpp16

type StopTransactionResponseJsonIdTagInfoStatus string

type StopTransactionResponseJsonIdTagInfo struct {
	// ExpiryDate corresponds to the JSON schema field "expiryDate".
	ExpiryDate *string `json:"expiryDate,omitempty" yaml:"expiryDate,omitempty" mapstructure:"expiryDate,omitempty"`

	// ParentIdTag corresponds to the JSON schema field "parentIdTag".
	ParentIdTag *string `json:"parentIdTag,omitempty" yaml:"parentIdTag,omitempty" mapstructure:"parentIdTag,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status StopTransactionResponseJsonIdTagInfoStatus `json:"status" yaml:"status" mapstructure:"status"`
}

const StopTransactionResponseJsonIdTagInfoStatusAccepted StopTransactionResponseJsonIdTagInfoStatus = "Accepted"
const StopTransactionResponseJsonIdTagInfoStatusBlocked StopTransactionResponseJsonIdTagInfoStatus = "Blocked"
const StopTransactionResponseJsonIdTagInfoStatusConcurrentTx StopTransactionResponseJsonIdTagInfoStatus = "ConcurrentTx"
const StopTransactionResponseJsonIdTagInfoStatusExpired StopTransactionResponseJsonIdTagInfoStatus = "Expired"
const StopTransactionResponseJsonIdTagInfoStatusInvalid StopTransactionResponseJsonIdTagInfoStatus = "Invalid"

type StopTransactionResponseJson struct {
	// IdTagInfo corresponds to the JSON schema field "idTagInfo".
	IdTagInfo *StopTransactionResponseJsonIdTagInfo `json:"idTagInfo,omitempty" yaml:"idTagInfo,omitempty" mapstructure:"idTagInfo,omitempty"`
}

func (*StopTransactionResponseJson) IsResponse() {}
