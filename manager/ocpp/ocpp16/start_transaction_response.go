// SPDX-License-Identifier: Apache-2.0

package ocpp16

type StartTransactionResponseJsonIdTagInfoStatus string

type StartTransactionResponseJsonIdTagInfo struct {
	// ExpiryDate corresponds to the JSON schema field "expiryDate".
	ExpiryDate *string `json:"expiryDate,omitempty" yaml:"expiryDate,omitempty" mapstructure:"expiryDate,omitempty"`

	// ParentIdTag corresponds to the JSON schema field "parentIdTag".
	ParentIdTag *string `json:"parentIdTag,omitempty" yaml:"parentIdTag,omitempty" mapstructure:"parentIdTag,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status StartTransactionResponseJsonIdTagInfoStatus `json:"status" yaml:"status" mapstructure:"status"`
}

const StartTransactionResponseJsonIdTagInfoStatusAccepted StartTransactionResponseJsonIdTagInfoStatus = "Accepted"
const StartTransactionResponseJsonIdTagInfoStatusBlocked StartTransactionResponseJsonIdTagInfoStatus = "Blocked"
const StartTransactionResponseJsonIdTagInfoStatusConcurrentTx StartTransactionResponseJsonIdTagInfoStatus = "ConcurrentTx"
const StartTransactionResponseJsonIdTagInfoStatusExpired StartTransactionResponseJsonIdTagInfoStatus = "Expired"
const StartTransactionResponseJsonIdTagInfoStatusInvalid StartTransactionResponseJsonIdTagInfoStatus = "Invalid"

type StartTransactionResponseJson struct {
	// IdTagInfo corresponds to the JSON schema field "idTagInfo".
	IdTagInfo StartTransactionResponseJsonIdTagInfo `json:"idTagInfo" yaml:"idTagInfo" mapstructure:"idTagInfo"`

	// TransactionId corresponds to the JSON schema field "transactionId".
	TransactionId int `json:"transactionId" yaml:"transactionId" mapstructure:"transactionId"`
}

func (*StartTransactionResponseJson) IsResponse() {}
