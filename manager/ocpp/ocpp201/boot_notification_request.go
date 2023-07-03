package ocpp201

type BootNotificationRequestJson struct {
	// ChargingStation corresponds to the JSON schema field "chargingStation".
	ChargingStation ChargingStationType `json:"chargingStation" yaml:"chargingStation" mapstructure:"chargingStation"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Reason corresponds to the JSON schema field "reason".
	Reason BootReasonEnumType `json:"reason" yaml:"reason" mapstructure:"reason"`
}

func (*BootNotificationRequestJson) IsRequest() {}

type BootReasonEnumType string

const BootReasonEnumTypeApplicationReset BootReasonEnumType = "ApplicationReset"
const BootReasonEnumTypeFirmwareUpdate BootReasonEnumType = "FirmwareUpdate"
const BootReasonEnumTypeLocalReset BootReasonEnumType = "LocalReset"
const BootReasonEnumTypePowerUp BootReasonEnumType = "PowerUp"
const BootReasonEnumTypeRemoteReset BootReasonEnumType = "RemoteReset"
const BootReasonEnumTypeScheduledReset BootReasonEnumType = "ScheduledReset"
const BootReasonEnumTypeTriggered BootReasonEnumType = "Triggered"
const BootReasonEnumTypeUnknown BootReasonEnumType = "Unknown"
const BootReasonEnumTypeWatchdog BootReasonEnumType = "Watchdog"

// Charge_ Point
// urn:x-oca:ocpp:uid:2:233122
// The physical system where an Electrical Vehicle (EV) can be charged.
type ChargingStationType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// This contains the firmware version of the Charging Station.
	//
	//
	FirmwareVersion *string `json:"firmwareVersion,omitempty" yaml:"firmwareVersion,omitempty" mapstructure:"firmwareVersion,omitempty"`

	// Device. Model. CI20_ Text
	// urn:x-oca:ocpp:uid:1:569325
	// Defines the model of the device.
	//
	Model string `json:"model" yaml:"model" mapstructure:"model"`

	// Modem corresponds to the JSON schema field "modem".
	Modem *ModemType `json:"modem,omitempty" yaml:"modem,omitempty" mapstructure:"modem,omitempty"`

	// Device. Serial_ Number. Serial_ Number
	// urn:x-oca:ocpp:uid:1:569324
	// Vendor-specific device identifier.
	//
	SerialNumber *string `json:"serialNumber,omitempty" yaml:"serialNumber,omitempty" mapstructure:"serialNumber,omitempty"`

	// Identifies the vendor (not necessarily in a unique manner).
	//
	VendorName string `json:"vendorName" yaml:"vendorName" mapstructure:"vendorName"`
}

// Wireless_ Communication_ Module
// urn:x-oca:ocpp:uid:2:233306
// Defines parameters required for initiating and maintaining wireless
// communication with other devices.
type ModemType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Wireless_ Communication_ Module. ICCID. CI20_ Text
	// urn:x-oca:ocpp:uid:1:569327
	// This contains the ICCID of the modem’s SIM card.
	//
	Iccid *string `json:"iccid,omitempty" yaml:"iccid,omitempty" mapstructure:"iccid,omitempty"`

	// Wireless_ Communication_ Module. IMSI. CI20_ Text
	// urn:x-oca:ocpp:uid:1:569328
	// This contains the IMSI of the modem’s SIM card.
	//
	Imsi *string `json:"imsi,omitempty" yaml:"imsi,omitempty" mapstructure:"imsi,omitempty"`
}
