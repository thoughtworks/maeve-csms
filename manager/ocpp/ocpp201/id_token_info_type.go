// SPDX-License-Identifier: Apache-2.0

package ocpp201

// ID_ Token
// urn:x-oca:ocpp:uid:2:233247
// Contains status information about an identifier.
// It is advised to not stop charging for a token that expires during charging, as
// ExpiryDate is only used for caching purposes. If ExpiryDate is not given, the
// status has no end date.
type IdTokenInfoType struct {
	// ID_ Token. Expiry. Date_ Time
	// urn:x-oca:ocpp:uid:1:569373
	// Date and Time after which the token must be considered invalid.
	//
	CacheExpiryDateTime *string `json:"cacheExpiryDateTime,omitempty" yaml:"cacheExpiryDateTime,omitempty" mapstructure:"cacheExpiryDateTime,omitempty"`

	// Priority from a business point of view. Default priority is 0, The range is
	// from -9 to 9. Higher values indicate a higher priority. The chargingPriority in
	// &lt;&lt;transactioneventresponse,TransactionEventResponse&gt;&gt; overrules
	// this one.
	//
	ChargingPriority *int `json:"chargingPriority,omitempty" yaml:"chargingPriority,omitempty" mapstructure:"chargingPriority,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Only used when the IdToken is only valid for one or more specific EVSEs, not
	// for the entire Charging Station.
	//
	//
	EvseId []int `json:"evseId,omitempty" yaml:"evseId,omitempty" mapstructure:"evseId,omitempty"`

	// GroupIdToken corresponds to the JSON schema field "groupIdToken".
	GroupIdToken *IdTokenType `json:"groupIdToken,omitempty" yaml:"groupIdToken,omitempty" mapstructure:"groupIdToken,omitempty"`

	// ID_ Token. Language1. Language_ Code
	// urn:x-oca:ocpp:uid:1:569374
	// Preferred user interface language of identifier user. Contains a language code
	// as defined in &lt;&lt;ref-RFC5646,[RFC5646]&gt;&gt;.
	//
	//
	Language1 *string `json:"language1,omitempty" yaml:"language1,omitempty" mapstructure:"language1,omitempty"`

	// ID_ Token. Language2. Language_ Code
	// urn:x-oca:ocpp:uid:1:569375
	// Second preferred user interface language of identifier user. Don’t use when
	// language1 is omitted, has to be different from language1. Contains a language
	// code as defined in &lt;&lt;ref-RFC5646,[RFC5646]&gt;&gt;.
	//
	Language2 *string `json:"language2,omitempty" yaml:"language2,omitempty" mapstructure:"language2,omitempty"`

	// PersonalMessage corresponds to the JSON schema field "personalMessage".
	PersonalMessage *MessageContentType `json:"personalMessage,omitempty" yaml:"personalMessage,omitempty" mapstructure:"personalMessage,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status AuthorizationStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`
}
