// SPDX-License-Identifier: Apache-2.0

package ocpp16

type ChangeConfigurationJson struct {
	// Key corresponds to the JSON schema field "key".
	Key string `json:"key" yaml:"key" mapstructure:"key"`

	// Value corresponds to the JSON schema field "value".
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

func (*ChangeConfigurationJson) IsRequest() {}
