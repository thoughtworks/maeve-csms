// SPDX-License-Identifier: Apache-2.0

package ocpp201

// Contains the identifier to use for authorization.
type AuthorizationData struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// IdToken corresponds to the JSON schema field "idToken".
	IdToken IdTokenType `json:"idToken" yaml:"idToken" mapstructure:"idToken"`

	// IdTokenInfo corresponds to the JSON schema field "idTokenInfo".
	IdTokenInfo *IdTokenInfoType `json:"idTokenInfo,omitempty" yaml:"idTokenInfo,omitempty" mapstructure:"idTokenInfo,omitempty"`
}

type UpdateEnumType string

const UpdateEnumTypeDifferential UpdateEnumType = "Differential"
const UpdateEnumTypeFull UpdateEnumType = "Full"

type SendLocalListRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// LocalAuthorizationList corresponds to the JSON schema field
	// "localAuthorizationList".
	LocalAuthorizationList []AuthorizationData `json:"localAuthorizationList,omitempty" yaml:"localAuthorizationList,omitempty" mapstructure:"localAuthorizationList,omitempty"`

	// UpdateType corresponds to the JSON schema field "updateType".
	UpdateType UpdateEnumType `json:"updateType" yaml:"updateType" mapstructure:"updateType"`

	// In case of a full update this is the version number of the full list. In case
	// of a differential update it is the version number of the list after the update
	// has been applied.
	//
	VersionNumber int `json:"versionNumber" yaml:"versionNumber" mapstructure:"versionNumber"`
}

func (*SendLocalListRequestJson) IsRequest() {}
