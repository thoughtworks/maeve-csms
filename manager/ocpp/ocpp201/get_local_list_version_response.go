// SPDX-License-Identifier: Apache-2.0

package ocpp201

type GetLocalListVersionResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// This contains the current version number of the local authorization list in the
	// Charging Station.
	//
	VersionNumber int `json:"versionNumber" yaml:"versionNumber" mapstructure:"versionNumber"`
}

func (*GetLocalListVersionResponseJson) IsResponse() {}
