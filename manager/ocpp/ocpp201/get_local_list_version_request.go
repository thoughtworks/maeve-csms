// SPDX-License-Identifier: Apache-2.0

package ocpp201

type GetLocalListVersionRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

func (*GetLocalListVersionRequestJson) IsRequest() {}
