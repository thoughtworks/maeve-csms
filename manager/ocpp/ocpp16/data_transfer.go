// SPDX-License-Identifier: Apache-2.0

package ocpp16

type DataTransferJson struct {
	// Data corresponds to the JSON schema field "data".
	Data *string `json:"data,omitempty" yaml:"data,omitempty" mapstructure:"data,omitempty"`

	// MessageId corresponds to the JSON schema field "messageId".
	MessageId *string `json:"messageId,omitempty" yaml:"messageId,omitempty" mapstructure:"messageId,omitempty"`

	// VendorId corresponds to the JSON schema field "vendorId".
	VendorId string `json:"vendorId" yaml:"vendorId" mapstructure:"vendorId"`
}

func (*DataTransferJson) IsRequest() {}
