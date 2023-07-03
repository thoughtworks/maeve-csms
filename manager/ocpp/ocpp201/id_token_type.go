package ocpp201

// Contains a case insensitive identifier to use for the authorization and the type
// of authorization to support multiple forms of identifiers.
type IdTokenType struct {
	// AdditionalInfo corresponds to the JSON schema field "additionalInfo".
	AdditionalInfo *[]AdditionalInfoType `json:"additionalInfo,omitempty" yaml:"additionalInfo,omitempty" mapstructure:"additionalInfo,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// IdToken is case insensitive. Might hold the hidden id of an RFID tag, but can
	// for example also contain a UUID.
	//
	IdToken string `json:"idToken" yaml:"idToken" mapstructure:"idToken"`

	// Type corresponds to the JSON schema field "type".
	Type IdTokenEnumType `json:"type" yaml:"type" mapstructure:"type"`
}
