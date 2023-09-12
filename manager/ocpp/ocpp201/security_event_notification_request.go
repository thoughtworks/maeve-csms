// SPDX-License-Identifier: Apache-2.0

package ocpp201

type SecurityEventNotificationRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Additional information about the occurred security event.
	//
	TechInfo *string `json:"techInfo,omitempty" yaml:"techInfo,omitempty" mapstructure:"techInfo,omitempty"`

	// Date and time at which the event occurred.
	//
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`

	// Type of the security event. This value should be taken from the Security events
	// list.
	//
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

func (*SecurityEventNotificationRequestJson) IsRequest() {}
