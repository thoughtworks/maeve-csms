// SPDX-License-Identifier: Apache-2.0

package ocpp201

type RequestStartTransactionResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status RequestStartStopStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`

	// When the transaction was already started by the Charging Station before the
	// RequestStartTransactionRequest was received, for example: cable plugged in
	// first. This contains the transactionId of the already started transaction.
	//
	TransactionId *string `json:"transactionId,omitempty" yaml:"transactionId,omitempty" mapstructure:"transactionId,omitempty"`
}

func (*RequestStartTransactionResponseJson) IsResponse() {}
