package ocpp16

const RemoteStartTransactionJsonChargingProfileChargingScheduleChargingRateUnitA RemoteStartTransactionJsonChargingProfileChargingScheduleChargingRateUnit = "A"
const RemoteStartTransactionJsonChargingProfileChargingProfileKindAbsolute RemoteStartTransactionJsonChargingProfileChargingProfileKind = "Absolute"
const RemoteStartTransactionJsonChargingProfileChargingProfileKindRecurring RemoteStartTransactionJsonChargingProfileChargingProfileKind = "Recurring"
const RemoteStartTransactionJsonChargingProfileChargingProfileKindRelative RemoteStartTransactionJsonChargingProfileChargingProfileKind = "Relative"

type RemoteStartTransactionJsonChargingProfileChargingProfilePurpose string

const RemoteStartTransactionJsonChargingProfileChargingProfilePurposeChargePointMaxProfile RemoteStartTransactionJsonChargingProfileChargingProfilePurpose = "ChargePointMaxProfile"
const RemoteStartTransactionJsonChargingProfileChargingProfilePurposeTxDefaultProfile RemoteStartTransactionJsonChargingProfileChargingProfilePurpose = "TxDefaultProfile"
const RemoteStartTransactionJsonChargingProfileChargingProfilePurposeTxProfile RemoteStartTransactionJsonChargingProfileChargingProfilePurpose = "TxProfile"

type RemoteStartTransactionJsonChargingProfileChargingScheduleChargingRateUnit string

const RemoteStartTransactionJsonChargingProfileChargingScheduleChargingRateUnitW RemoteStartTransactionJsonChargingProfileChargingScheduleChargingRateUnit = "W"

type RemoteStartTransactionJsonChargingProfileChargingScheduleChargingSchedulePeriodElem struct {
	// Limit corresponds to the JSON schema field "limit".
	Limit float64 `json:"limit" yaml:"limit" mapstructure:"limit"`

	// NumberPhases corresponds to the JSON schema field "numberPhases".
	NumberPhases *int `json:"numberPhases,omitempty" yaml:"numberPhases,omitempty" mapstructure:"numberPhases,omitempty"`

	// StartPeriod corresponds to the JSON schema field "startPeriod".
	StartPeriod int `json:"startPeriod" yaml:"startPeriod" mapstructure:"startPeriod"`
}

type RemoteStartTransactionJsonChargingProfileChargingSchedule struct {
	// ChargingRateUnit corresponds to the JSON schema field "chargingRateUnit".
	ChargingRateUnit RemoteStartTransactionJsonChargingProfileChargingScheduleChargingRateUnit `json:"chargingRateUnit" yaml:"chargingRateUnit" mapstructure:"chargingRateUnit"`

	// ChargingSchedulePeriod corresponds to the JSON schema field
	// "chargingSchedulePeriod".
	ChargingSchedulePeriod []RemoteStartTransactionJsonChargingProfileChargingScheduleChargingSchedulePeriodElem `json:"chargingSchedulePeriod" yaml:"chargingSchedulePeriod" mapstructure:"chargingSchedulePeriod"`

	// Duration corresponds to the JSON schema field "duration".
	Duration *int `json:"duration,omitempty" yaml:"duration,omitempty" mapstructure:"duration,omitempty"`

	// MinChargingRate corresponds to the JSON schema field "minChargingRate".
	MinChargingRate *float64 `json:"minChargingRate,omitempty" yaml:"minChargingRate,omitempty" mapstructure:"minChargingRate,omitempty"`

	// StartSchedule corresponds to the JSON schema field "startSchedule".
	StartSchedule *string `json:"startSchedule,omitempty" yaml:"startSchedule,omitempty" mapstructure:"startSchedule,omitempty"`
}

type RemoteStartTransactionJsonChargingProfileRecurrencyKind string

const RemoteStartTransactionJsonChargingProfileRecurrencyKindDaily RemoteStartTransactionJsonChargingProfileRecurrencyKind = "Daily"
const RemoteStartTransactionJsonChargingProfileRecurrencyKindWeekly RemoteStartTransactionJsonChargingProfileRecurrencyKind = "Weekly"

type RemoteStartTransactionJsonChargingProfile struct {
	// ChargingProfileId corresponds to the JSON schema field "chargingProfileId".
	ChargingProfileId int `json:"chargingProfileId" yaml:"chargingProfileId" mapstructure:"chargingProfileId"`

	// ChargingProfileKind corresponds to the JSON schema field "chargingProfileKind".
	ChargingProfileKind RemoteStartTransactionJsonChargingProfileChargingProfileKind `json:"chargingProfileKind" yaml:"chargingProfileKind" mapstructure:"chargingProfileKind"`

	// ChargingProfilePurpose corresponds to the JSON schema field
	// "chargingProfilePurpose".
	ChargingProfilePurpose RemoteStartTransactionJsonChargingProfileChargingProfilePurpose `json:"chargingProfilePurpose" yaml:"chargingProfilePurpose" mapstructure:"chargingProfilePurpose"`

	// ChargingSchedule corresponds to the JSON schema field "chargingSchedule".
	ChargingSchedule RemoteStartTransactionJsonChargingProfileChargingSchedule `json:"chargingSchedule" yaml:"chargingSchedule" mapstructure:"chargingSchedule"`

	// RecurrencyKind corresponds to the JSON schema field "recurrencyKind".
	RecurrencyKind *RemoteStartTransactionJsonChargingProfileRecurrencyKind `json:"recurrencyKind,omitempty" yaml:"recurrencyKind,omitempty" mapstructure:"recurrencyKind,omitempty"`

	// StackLevel corresponds to the JSON schema field "stackLevel".
	StackLevel int `json:"stackLevel" yaml:"stackLevel" mapstructure:"stackLevel"`

	// TransactionId corresponds to the JSON schema field "transactionId".
	TransactionId *int `json:"transactionId,omitempty" yaml:"transactionId,omitempty" mapstructure:"transactionId,omitempty"`

	// ValidFrom corresponds to the JSON schema field "validFrom".
	ValidFrom *string `json:"validFrom,omitempty" yaml:"validFrom,omitempty" mapstructure:"validFrom,omitempty"`

	// ValidTo corresponds to the JSON schema field "validTo".
	ValidTo *string `json:"validTo,omitempty" yaml:"validTo,omitempty" mapstructure:"validTo,omitempty"`
}

type RemoteStartTransactionJson struct {
	// ChargingProfile corresponds to the JSON schema field "chargingProfile".
	ChargingProfile *RemoteStartTransactionJsonChargingProfile `json:"chargingProfile,omitempty" yaml:"chargingProfile,omitempty" mapstructure:"chargingProfile,omitempty"`

	// ConnectorId corresponds to the JSON schema field "connectorId".
	ConnectorId *int `json:"connectorId,omitempty" yaml:"connectorId,omitempty" mapstructure:"connectorId,omitempty"`

	// IdTag corresponds to the JSON schema field "idTag".
	IdTag string `json:"idTag" yaml:"idTag" mapstructure:"idTag"`
}

type RemoteStartTransactionJsonChargingProfileChargingProfileKind string

func (*RemoteStartTransactionJson) IsRequest() {}
