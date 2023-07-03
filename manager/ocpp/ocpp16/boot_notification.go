package ocpp16

type BootNotificationJson struct {
	// ChargeBoxSerialNumber corresponds to the JSON schema field
	// "chargeBoxSerialNumber".
	ChargeBoxSerialNumber *string `json:"chargeBoxSerialNumber,omitempty" yaml:"chargeBoxSerialNumber,omitempty" mapstructure:"chargeBoxSerialNumber,omitempty"`

	// ChargePointModel corresponds to the JSON schema field "chargePointModel".
	ChargePointModel string `json:"chargePointModel" yaml:"chargePointModel" mapstructure:"chargePointModel"`

	// ChargePointSerialNumber corresponds to the JSON schema field
	// "chargePointSerialNumber".
	ChargePointSerialNumber *string `json:"chargePointSerialNumber,omitempty" yaml:"chargePointSerialNumber,omitempty" mapstructure:"chargePointSerialNumber,omitempty"`

	// ChargePointVendor corresponds to the JSON schema field "chargePointVendor".
	ChargePointVendor string `json:"chargePointVendor" yaml:"chargePointVendor" mapstructure:"chargePointVendor"`

	// FirmwareVersion corresponds to the JSON schema field "firmwareVersion".
	FirmwareVersion *string `json:"firmwareVersion,omitempty" yaml:"firmwareVersion,omitempty" mapstructure:"firmwareVersion,omitempty"`

	// Iccid corresponds to the JSON schema field "iccid".
	Iccid *string `json:"iccid,omitempty" yaml:"iccid,omitempty" mapstructure:"iccid,omitempty"`

	// Imsi corresponds to the JSON schema field "imsi".
	Imsi *string `json:"imsi,omitempty" yaml:"imsi,omitempty" mapstructure:"imsi,omitempty"`

	// MeterSerialNumber corresponds to the JSON schema field "meterSerialNumber".
	MeterSerialNumber *string `json:"meterSerialNumber,omitempty" yaml:"meterSerialNumber,omitempty" mapstructure:"meterSerialNumber,omitempty"`

	// MeterType corresponds to the JSON schema field "meterType".
	MeterType *string `json:"meterType,omitempty" yaml:"meterType,omitempty" mapstructure:"meterType,omitempty"`
}

func (*BootNotificationJson) IsRequest() {}
