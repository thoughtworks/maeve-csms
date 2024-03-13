// SPDX-License-Identifier: Apache-2.0

package ocpp201

type RequestStopTransactionRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// The identifier of the transaction which the Charging Station is requested to
	// stop.
	//
	TransactionId string `json:"transactionId" yaml:"transactionId" mapstructure:"transactionId"`
}

func (*RequestStopTransactionRequestJson) IsRequest() {}
