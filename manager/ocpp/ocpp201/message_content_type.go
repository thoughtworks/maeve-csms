// SPDX-License-Identifier: Apache-2.0

package ocpp201

// Message_ Content
// urn:x-enexis:ecdm:uid:2:234490
// Contains message details, for a message to be displayed on a Charging Station.
type MessageContentType struct {
	// Message_ Content. Content. Message
	// urn:x-enexis:ecdm:uid:1:570852
	// Message contents.
	//
	//
	Content string `json:"content" yaml:"content" mapstructure:"content"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Format corresponds to the JSON schema field "format".
	Format MessageFormatEnumType `json:"format" yaml:"format" mapstructure:"format"`

	// Message_ Content. Language. Language_ Code
	// urn:x-enexis:ecdm:uid:1:570849
	// Message language identifier. Contains a language code as defined in
	// &lt;&lt;ref-RFC5646,[RFC5646]&gt;&gt;.
	//
	Language *string `json:"language,omitempty" yaml:"language,omitempty" mapstructure:"language,omitempty"`
}

const MessageFormatEnumTypeASCII MessageFormatEnumType = "ASCII"
const MessageFormatEnumTypeHTML MessageFormatEnumType = "HTML"
const MessageFormatEnumTypeURI MessageFormatEnumType = "URI"
const MessageFormatEnumTypeUTF8 MessageFormatEnumType = "UTF8"
