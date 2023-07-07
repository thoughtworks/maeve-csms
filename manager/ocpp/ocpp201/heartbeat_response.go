// SPDX-License-Identifier: Apache-2.0

package ocpp201

type HeartbeatResponseJson struct {
	// Contains the current time of the CSMS.
	//
	CurrentTime string `json:"currentTime" yaml:"currentTime" mapstructure:"currentTime"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

func (*HeartbeatResponseJson) IsResponse() {}
