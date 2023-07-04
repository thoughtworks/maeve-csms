package ocpp16

type StatusNotificationJson struct {
	// ConnectorId corresponds to the JSON schema field "connectorId".
	ConnectorId int `json:"connectorId" yaml:"connectorId" mapstructure:"connectorId"`

	// ErrorCode corresponds to the JSON schema field "errorCode".
	ErrorCode StatusNotificationJsonErrorCode `json:"errorCode" yaml:"errorCode" mapstructure:"errorCode"`

	// Info corresponds to the JSON schema field "info".
	Info *string `json:"info,omitempty" yaml:"info,omitempty" mapstructure:"info,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status StatusNotificationJsonStatus `json:"status" yaml:"status" mapstructure:"status"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp *string `json:"timestamp,omitempty" yaml:"timestamp,omitempty" mapstructure:"timestamp,omitempty"`

	// VendorErrorCode corresponds to the JSON schema field "vendorErrorCode".
	VendorErrorCode *string `json:"vendorErrorCode,omitempty" yaml:"vendorErrorCode,omitempty" mapstructure:"vendorErrorCode,omitempty"`

	// VendorId corresponds to the JSON schema field "vendorId".
	VendorId *string `json:"vendorId,omitempty" yaml:"vendorId,omitempty" mapstructure:"vendorId,omitempty"`
}

func (*StatusNotificationJson) IsRequest() {}

type StatusNotificationJsonErrorCode string

const StatusNotificationJsonErrorCodeConnectorLockFailure StatusNotificationJsonErrorCode = "ConnectorLockFailure"
const StatusNotificationJsonErrorCodeEVCommunicationError StatusNotificationJsonErrorCode = "EVCommunicationError"
const StatusNotificationJsonErrorCodeGroundFailure StatusNotificationJsonErrorCode = "GroundFailure"
const StatusNotificationJsonErrorCodeHighTemperature StatusNotificationJsonErrorCode = "HighTemperature"
const StatusNotificationJsonErrorCodeInternalError StatusNotificationJsonErrorCode = "InternalError"
const StatusNotificationJsonErrorCodeLocalListConflict StatusNotificationJsonErrorCode = "LocalListConflict"
const StatusNotificationJsonErrorCodeNoError StatusNotificationJsonErrorCode = "NoError"
const StatusNotificationJsonErrorCodeOtherError StatusNotificationJsonErrorCode = "OtherError"
const StatusNotificationJsonErrorCodeOverCurrentFailure StatusNotificationJsonErrorCode = "OverCurrentFailure"
const StatusNotificationJsonErrorCodeOverVoltage StatusNotificationJsonErrorCode = "OverVoltage"
const StatusNotificationJsonErrorCodePowerMeterFailure StatusNotificationJsonErrorCode = "PowerMeterFailure"
const StatusNotificationJsonErrorCodePowerSwitchFailure StatusNotificationJsonErrorCode = "PowerSwitchFailure"
const StatusNotificationJsonErrorCodeReaderFailure StatusNotificationJsonErrorCode = "ReaderFailure"
const StatusNotificationJsonErrorCodeResetFailure StatusNotificationJsonErrorCode = "ResetFailure"
const StatusNotificationJsonErrorCodeUnderVoltage StatusNotificationJsonErrorCode = "UnderVoltage"
const StatusNotificationJsonErrorCodeWeakSignal StatusNotificationJsonErrorCode = "WeakSignal"

type StatusNotificationJsonStatus string

const StatusNotificationJsonStatusAvailable StatusNotificationJsonStatus = "Available"
const StatusNotificationJsonStatusCharging StatusNotificationJsonStatus = "Charging"
const StatusNotificationJsonStatusFaulted StatusNotificationJsonStatus = "Faulted"
const StatusNotificationJsonStatusFinishing StatusNotificationJsonStatus = "Finishing"
const StatusNotificationJsonStatusPreparing StatusNotificationJsonStatus = "Preparing"
const StatusNotificationJsonStatusReserved StatusNotificationJsonStatus = "Reserved"
const StatusNotificationJsonStatusSuspendedEV StatusNotificationJsonStatus = "SuspendedEV"
const StatusNotificationJsonStatusSuspendedEVSE StatusNotificationJsonStatus = "SuspendedEVSE"
const StatusNotificationJsonStatusUnavailable StatusNotificationJsonStatus = "Unavailable"
