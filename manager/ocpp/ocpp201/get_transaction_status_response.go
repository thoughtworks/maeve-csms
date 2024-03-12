// SPDX-License-Identifier: Apache-2.0

package ocpp201

type GetTransactionStatusResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Whether there are still message to be delivered.
	//
	MessagesInQueue bool `json:"messagesInQueue" yaml:"messagesInQueue" mapstructure:"messagesInQueue"`

	// Whether the transaction is still ongoing.
	//
	OngoingIndicator *bool `json:"ongoingIndicator,omitempty" yaml:"ongoingIndicator,omitempty" mapstructure:"ongoingIndicator,omitempty"`
}

func (*GetTransactionStatusResponseJson) IsResponse() {}
