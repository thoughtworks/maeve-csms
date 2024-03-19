// SPDX-License-Identifier: Apache-2.0

package ocpp201

type APNAuthenticationEnumType string

const APNAuthenticationEnumTypeAUTO APNAuthenticationEnumType = "AUTO"
const APNAuthenticationEnumTypeCHAP APNAuthenticationEnumType = "CHAP"
const APNAuthenticationEnumTypeNONE APNAuthenticationEnumType = "NONE"
const APNAuthenticationEnumTypePAP APNAuthenticationEnumType = "PAP"

// APN
// urn:x-oca:ocpp:uid:2:233134
// Collection of configuration data needed to make a data-connection over a
// cellular network.
//
// NOTE: When asking a GSM modem to dial in, it is possible to specify which mobile
// operator should be used. This can be done with the mobile country code (MCC) in
// combination with a mobile network code (MNC). Example: If your preferred network
// is Vodafone Netherlands, the MCC=204 and the MNC=04 which means the key
// PreferredNetwork = 20404 Some modems allows to specify a preferred network,
// which means, if this network is not available, a different network is used. If
// you specify UseOnlyPreferredNetwork and this network is not available, the modem
// will not dial in.
type APNType struct {
	// APN. APN. URI
	// urn:x-oca:ocpp:uid:1:568814
	// The Access Point Name as an URL.
	//
	Apn string `json:"apn" yaml:"apn" mapstructure:"apn"`

	// ApnAuthentication corresponds to the JSON schema field "apnAuthentication".
	ApnAuthentication APNAuthenticationEnumType `json:"apnAuthentication" yaml:"apnAuthentication" mapstructure:"apnAuthentication"`

	// APN. APN. Password
	// urn:x-oca:ocpp:uid:1:568819
	// APN Password.
	//
	ApnPassword *string `json:"apnPassword,omitempty" yaml:"apnPassword,omitempty" mapstructure:"apnPassword,omitempty"`

	// APN. APN. User_ Name
	// urn:x-oca:ocpp:uid:1:568818
	// APN username.
	//
	ApnUserName *string `json:"apnUserName,omitempty" yaml:"apnUserName,omitempty" mapstructure:"apnUserName,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// APN. Preferred_ Network. Mobile_ Network_ ID
	// urn:x-oca:ocpp:uid:1:568822
	// Preferred network, written as MCC and MNC concatenated. See note.
	//
	PreferredNetwork *string `json:"preferredNetwork,omitempty" yaml:"preferredNetwork,omitempty" mapstructure:"preferredNetwork,omitempty"`

	// APN. SIMPIN. PIN_ Code
	// urn:x-oca:ocpp:uid:1:568821
	// SIM card pin code.
	//
	SimPin *int `json:"simPin,omitempty" yaml:"simPin,omitempty" mapstructure:"simPin,omitempty"`

	// APN. Use_ Only_ Preferred_ Network. Indicator
	// urn:x-oca:ocpp:uid:1:568824
	// Default: false. Use only the preferred Network, do
	// not dial in when not available. See Note.
	//
	UseOnlyPreferredNetwork bool `json:"useOnlyPreferredNetwork,omitempty" yaml:"useOnlyPreferredNetwork,omitempty" mapstructure:"useOnlyPreferredNetwork,omitempty"`
}

// Communication_ Function
// urn:x-oca:ocpp:uid:2:233304
// The NetworkConnectionProfile defines the functional and technical parameters of
// a communication link.
type NetworkConnectionProfileType struct {
	// Apn corresponds to the JSON schema field "apn".
	Apn *APNType `json:"apn,omitempty" yaml:"apn,omitempty" mapstructure:"apn,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Duration in seconds before a message send by the Charging Station via this
	// network connection times-out.
	// The best setting depends on the underlying network and response times of the
	// CSMS.
	// If you are looking for a some guideline: use 30 seconds as a starting point.
	//
	MessageTimeout int `json:"messageTimeout" yaml:"messageTimeout" mapstructure:"messageTimeout"`

	// Communication_ Function. OCPP_ Central_ System_ URL. URI
	// urn:x-oca:ocpp:uid:1:569357
	// URL of the CSMS(s) that this Charging Station  communicates with.
	//
	OcppCsmsUrl string `json:"ocppCsmsUrl" yaml:"ocppCsmsUrl" mapstructure:"ocppCsmsUrl"`

	// OcppInterface corresponds to the JSON schema field "ocppInterface".
	OcppInterface OCPPInterfaceEnumType `json:"ocppInterface" yaml:"ocppInterface" mapstructure:"ocppInterface"`

	// OcppTransport corresponds to the JSON schema field "ocppTransport".
	OcppTransport OCPPTransportEnumType `json:"ocppTransport" yaml:"ocppTransport" mapstructure:"ocppTransport"`

	// OcppVersion corresponds to the JSON schema field "ocppVersion".
	OcppVersion OCPPVersionEnumType `json:"ocppVersion" yaml:"ocppVersion" mapstructure:"ocppVersion"`

	// This field specifies the security profile used when connecting to the CSMS with
	// this NetworkConnectionProfile.
	//
	SecurityProfile int `json:"securityProfile" yaml:"securityProfile" mapstructure:"securityProfile"`

	// Vpn corresponds to the JSON schema field "vpn".
	Vpn *VPNType `json:"vpn,omitempty" yaml:"vpn,omitempty" mapstructure:"vpn,omitempty"`
}

type OCPPInterfaceEnumType string

const OCPPInterfaceEnumTypeWired0 OCPPInterfaceEnumType = "Wired0"
const OCPPInterfaceEnumTypeWired1 OCPPInterfaceEnumType = "Wired1"
const OCPPInterfaceEnumTypeWired2 OCPPInterfaceEnumType = "Wired2"
const OCPPInterfaceEnumTypeWired3 OCPPInterfaceEnumType = "Wired3"
const OCPPInterfaceEnumTypeWireless0 OCPPInterfaceEnumType = "Wireless0"
const OCPPInterfaceEnumTypeWireless1 OCPPInterfaceEnumType = "Wireless1"
const OCPPInterfaceEnumTypeWireless2 OCPPInterfaceEnumType = "Wireless2"
const OCPPInterfaceEnumTypeWireless3 OCPPInterfaceEnumType = "Wireless3"

type OCPPTransportEnumType string

const OCPPTransportEnumTypeJSON OCPPTransportEnumType = "JSON"
const OCPPTransportEnumTypeSOAP OCPPTransportEnumType = "SOAP"

type OCPPVersionEnumType string

const OCPPVersionEnumTypeOCPP12 OCPPVersionEnumType = "OCPP12"
const OCPPVersionEnumTypeOCPP15 OCPPVersionEnumType = "OCPP15"
const OCPPVersionEnumTypeOCPP16 OCPPVersionEnumType = "OCPP16"
const OCPPVersionEnumTypeOCPP20 OCPPVersionEnumType = "OCPP20"

type SetNetworkProfileRequestJson struct {
	// Slot in which the configuration should be stored.
	//
	ConfigurationSlot int `json:"configurationSlot" yaml:"configurationSlot" mapstructure:"configurationSlot"`

	// ConnectionData corresponds to the JSON schema field "connectionData".
	ConnectionData NetworkConnectionProfileType `json:"connectionData" yaml:"connectionData" mapstructure:"connectionData"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

func (*SetNetworkProfileRequestJson) IsRequest() {}

type VPNEnumType string

const VPNEnumTypeIKEv2 VPNEnumType = "IKEv2"
const VPNEnumTypeIPSec VPNEnumType = "IPSec"
const VPNEnumTypeL2TP VPNEnumType = "L2TP"
const VPNEnumTypePPTP VPNEnumType = "PPTP"

// VPN
// urn:x-oca:ocpp:uid:2:233268
// VPN Configuration settings
type VPNType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// VPN. Group. Group_ Name
	// urn:x-oca:ocpp:uid:1:569274
	// VPN group.
	//
	Group *string `json:"group,omitempty" yaml:"group,omitempty" mapstructure:"group,omitempty"`

	// VPN. Key. VPN_ Key
	// urn:x-oca:ocpp:uid:1:569276
	// VPN shared secret.
	//
	Key string `json:"key" yaml:"key" mapstructure:"key"`

	// VPN. Password. Password
	// urn:x-oca:ocpp:uid:1:569275
	// VPN Password.
	//
	Password string `json:"password" yaml:"password" mapstructure:"password"`

	// VPN. Server. URI
	// urn:x-oca:ocpp:uid:1:569272
	// VPN Server Address
	//
	Server string `json:"server" yaml:"server" mapstructure:"server"`

	// Type corresponds to the JSON schema field "type".
	Type VPNEnumType `json:"type" yaml:"type" mapstructure:"type"`

	// VPN. User. User_ Name
	// urn:x-oca:ocpp:uid:1:569273
	// VPN User
	//
	User string `json:"user" yaml:"user" mapstructure:"user"`
}
