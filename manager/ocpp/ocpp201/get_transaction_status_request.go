// SPDX-License-Identifier: Apache-2.0

package ocpp201

type GetTransactionStatusRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// The Id of the transaction for which the status is requested.
	//
	TransactionId *string `json:"transactionId,omitempty" yaml:"transactionId,omitempty" mapstructure:"transactionId,omitempty"`
}

func (*GetTransactionStatusRequestJson) IsRequest() {}
