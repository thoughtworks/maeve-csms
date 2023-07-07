// SPDX-License-Identifier: Apache-2.0

package ocpp201

// Contains a case insensitive identifier to use for the authorization and the type
// of authorization to support multiple forms of identifiers.
type AdditionalInfoType struct {
	// This field specifies the additional IdToken.
	//
	AdditionalIdToken string `json:"additionalIdToken" yaml:"additionalIdToken" mapstructure:"additionalIdToken"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// This defines the type of the additionalIdToken. This is a custom type, so the
	// implementation needs to be agreed upon by all involved parties.
	//
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}
