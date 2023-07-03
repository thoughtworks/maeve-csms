package ocpp201

type RegistrationStatusEnumType string

const RegistrationStatusEnumTypeAccepted RegistrationStatusEnumType = "Accepted"
const RegistrationStatusEnumTypePending RegistrationStatusEnumType = "Pending"
const RegistrationStatusEnumTypeRejected RegistrationStatusEnumType = "Rejected"

type BootNotificationResponseJson struct {
	// This contains the CSMS’s current time.
	//
	CurrentTime string `json:"currentTime" yaml:"currentTime" mapstructure:"currentTime"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// When &lt;&lt;cmn_registrationstatusenumtype,Status&gt;&gt; is Accepted, this
	// contains the heartbeat interval in seconds. If the CSMS returns something other
	// than Accepted, the value of the interval field indicates the minimum wait time
	// before sending a next BootNotification request.
	//
	Interval int `json:"interval" yaml:"interval" mapstructure:"interval"`

	// Status corresponds to the JSON schema field "status".
	Status RegistrationStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`

	// StatusInfo corresponds to the JSON schema field "statusInfo".
	StatusInfo *StatusInfoType `json:"statusInfo,omitempty" yaml:"statusInfo,omitempty" mapstructure:"statusInfo,omitempty"`
}

func (*BootNotificationResponseJson) IsResponse() {}
