// SPDX-License-Identifier: Apache-2.0

package ocpp16

type SecurityEventNotificationJson struct {
	// TechInfo corresponds to the JSON schema field "techInfo".
	TechInfo *string `json:"techInfo,omitempty" yaml:"techInfo,omitempty" mapstructure:"techInfo,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`

	// Type corresponds to the JSON schema field "type".
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

func (*SecurityEventNotificationJson) IsRequest() {}
