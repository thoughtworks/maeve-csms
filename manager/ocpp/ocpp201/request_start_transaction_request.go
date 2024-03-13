// SPDX-License-Identifier: Apache-2.0

package ocpp201

type ChargingProfileKindEnumType string

const ChargingProfileKindEnumTypeAbsolute ChargingProfileKindEnumType = "Absolute"
const ChargingProfileKindEnumTypeRecurring ChargingProfileKindEnumType = "Recurring"
const ChargingProfileKindEnumTypeRelative ChargingProfileKindEnumType = "Relative"

type ChargingProfilePurposeEnumType string

const ChargingProfilePurposeEnumTypeChargingStationExternalConstraints ChargingProfilePurposeEnumType = "ChargingStationExternalConstraints"
const ChargingProfilePurposeEnumTypeChargingStationMaxProfile ChargingProfilePurposeEnumType = "ChargingStationMaxProfile"
const ChargingProfilePurposeEnumTypeTxDefaultProfile ChargingProfilePurposeEnumType = "TxDefaultProfile"
const ChargingProfilePurposeEnumTypeTxProfile ChargingProfilePurposeEnumType = "TxProfile"

type ChargingRateUnitEnumType string

const ChargingRateUnitEnumTypeA ChargingRateUnitEnumType = "A"
const ChargingRateUnitEnumTypeW ChargingRateUnitEnumType = "W"

// Charging_ Schedule_ Period
// urn:x-oca:ocpp:uid:2:233257
// Charging schedule period structure defines a time period in a charging schedule.
type ChargingSchedulePeriodType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Charging_ Schedule_ Period. Limit. Measure
	// urn:x-oca:ocpp:uid:1:569241
	// Charging rate limit during the schedule period, in the applicable
	// chargingRateUnit, for example in Amperes (A) or Watts (W). Accepts at most one
	// digit fraction (e.g. 8.1).
	//
	Limit float64 `json:"limit" yaml:"limit" mapstructure:"limit"`

	// Charging_ Schedule_ Period. Number_ Phases. Counter
	// urn:x-oca:ocpp:uid:1:569242
	// The number of phases that can be used for charging. If a number of phases is
	// needed, numberPhases=3 will be assumed unless another number is given.
	//
	NumberPhases *int `json:"numberPhases,omitempty" yaml:"numberPhases,omitempty" mapstructure:"numberPhases,omitempty"`

	// Values: 1..3, Used if numberPhases=1 and if the EVSE is capable of switching
	// the phase connected to the EV, i.e. ACPhaseSwitchingSupported is defined and
	// true. It’s not allowed unless both conditions above are true. If both
	// conditions are true, and phaseToUse is omitted, the Charging Station / EVSE
	// will make the selection on its own.
	//
	//
	PhaseToUse *int `json:"phaseToUse,omitempty" yaml:"phaseToUse,omitempty" mapstructure:"phaseToUse,omitempty"`

	// Charging_ Schedule_ Period. Start_ Period. Elapsed_ Time
	// urn:x-oca:ocpp:uid:1:569240
	// Start of the period, in seconds from the start of schedule. The value of
	// StartPeriod also defines the stop time of the previous period.
	//
	StartPeriod int `json:"startPeriod" yaml:"startPeriod" mapstructure:"startPeriod"`
}

type CostKindEnumType string

const CostKindEnumTypeCarbonDioxideEmission CostKindEnumType = "CarbonDioxideEmission"
const CostKindEnumTypeRelativePricePercentage CostKindEnumType = "RelativePricePercentage"
const CostKindEnumTypeRenewableGenerationPercentage CostKindEnumType = "RenewableGenerationPercentage"

// Cost
// urn:x-oca:ocpp:uid:2:233258
type CostType struct {
	// Cost. Amount. Amount
	// urn:x-oca:ocpp:uid:1:569244
	// The estimated or actual cost per kWh
	//
	Amount int `json:"amount" yaml:"amount" mapstructure:"amount"`

	// Cost. Amount_ Multiplier. Integer
	// urn:x-oca:ocpp:uid:1:569245
	// Values: -3..3, The amountMultiplier defines the exponent to base 10 (dec). The
	// final value is determined by: amount * 10 ^ amountMultiplier
	//
	AmountMultiplier *int `json:"amountMultiplier,omitempty" yaml:"amountMultiplier,omitempty" mapstructure:"amountMultiplier,omitempty"`

	// CostKind corresponds to the JSON schema field "costKind".
	CostKind CostKindEnumType `json:"costKind" yaml:"costKind" mapstructure:"costKind"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`
}

// Consumption_ Cost
// urn:x-oca:ocpp:uid:2:233259
type ConsumptionCostType struct {
	// Cost corresponds to the JSON schema field "cost".
	Cost []CostType `json:"cost" yaml:"cost" mapstructure:"cost"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Consumption_ Cost. Start_ Value. Numeric
	// urn:x-oca:ocpp:uid:1:569246
	// The lowest level of consumption that defines the starting point of this
	// consumption block. The block interval extends to the start of the next
	// interval.
	//
	StartValue float64 `json:"startValue" yaml:"startValue" mapstructure:"startValue"`
}

// Relative_ Timer_ Interval
// urn:x-oca:ocpp:uid:2:233270
type RelativeTimeIntervalType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Relative_ Timer_ Interval. Duration. Elapsed_ Time
	// urn:x-oca:ocpp:uid:1:569280
	// Duration of the interval, in seconds.
	//
	Duration *int `json:"duration,omitempty" yaml:"duration,omitempty" mapstructure:"duration,omitempty"`

	// Relative_ Timer_ Interval. Start. Elapsed_ Time
	// urn:x-oca:ocpp:uid:1:569279
	// Start of the interval, in seconds from NOW.
	//
	Start int `json:"start" yaml:"start" mapstructure:"start"`
}

// Sales_ Tariff_ Entry
// urn:x-oca:ocpp:uid:2:233271
type SalesTariffEntryType struct {
	// ConsumptionCost corresponds to the JSON schema field "consumptionCost".
	ConsumptionCost []ConsumptionCostType `json:"consumptionCost,omitempty" yaml:"consumptionCost,omitempty" mapstructure:"consumptionCost,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Sales_ Tariff_ Entry. E_ Price_ Level. Unsigned_ Integer
	// urn:x-oca:ocpp:uid:1:569281
	// Defines the price level of this SalesTariffEntry (referring to
	// NumEPriceLevels). Small values for the EPriceLevel represent a cheaper
	// TariffEntry. Large values for the EPriceLevel represent a more expensive
	// TariffEntry.
	//
	EPriceLevel *int `json:"ePriceLevel,omitempty" yaml:"ePriceLevel,omitempty" mapstructure:"ePriceLevel,omitempty"`

	// RelativeTimeInterval corresponds to the JSON schema field
	// "relativeTimeInterval".
	RelativeTimeInterval RelativeTimeIntervalType `json:"relativeTimeInterval" yaml:"relativeTimeInterval" mapstructure:"relativeTimeInterval"`
}

// Sales_ Tariff
// urn:x-oca:ocpp:uid:2:233272
// NOTE: This dataType is based on dataTypes from &lt;&lt;ref-ISOIEC15118-2,ISO
// 15118-2&gt;&gt;.
type SalesTariffType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Identified_ Object. MRID. Numeric_ Identifier
	// urn:x-enexis:ecdm:uid:1:569198
	// SalesTariff identifier used to identify one sales tariff. An SAID remains a
	// unique identifier for one schedule throughout a charging session.
	//
	Id int `json:"id" yaml:"id" mapstructure:"id"`

	// Sales_ Tariff. Num_ E_ Price_ Levels. Counter
	// urn:x-oca:ocpp:uid:1:569284
	// Defines the overall number of distinct price levels used across all provided
	// SalesTariff elements.
	//
	NumEPriceLevels *int `json:"numEPriceLevels,omitempty" yaml:"numEPriceLevels,omitempty" mapstructure:"numEPriceLevels,omitempty"`

	// Sales_ Tariff. Sales. Tariff_ Description
	// urn:x-oca:ocpp:uid:1:569283
	// A human readable title/short description of the sales tariff e.g. for HMI
	// display purposes.
	//
	SalesTariffDescription *string `json:"salesTariffDescription,omitempty" yaml:"salesTariffDescription,omitempty" mapstructure:"salesTariffDescription,omitempty"`

	// SalesTariffEntry corresponds to the JSON schema field "salesTariffEntry".
	SalesTariffEntry []SalesTariffEntryType `json:"salesTariffEntry" yaml:"salesTariffEntry" mapstructure:"salesTariffEntry"`
}

// Charging_ Schedule
// urn:x-oca:ocpp:uid:2:233256
// Charging schedule structure defines a list of charging periods, as used in:
// GetCompositeSchedule.conf and ChargingProfile.
type ChargingScheduleType struct {
	// ChargingRateUnit corresponds to the JSON schema field "chargingRateUnit".
	ChargingRateUnit ChargingRateUnitEnumType `json:"chargingRateUnit" yaml:"chargingRateUnit" mapstructure:"chargingRateUnit"`

	// ChargingSchedulePeriod corresponds to the JSON schema field
	// "chargingSchedulePeriod".
	ChargingSchedulePeriod []ChargingSchedulePeriodType `json:"chargingSchedulePeriod" yaml:"chargingSchedulePeriod" mapstructure:"chargingSchedulePeriod"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Charging_ Schedule. Duration. Elapsed_ Time
	// urn:x-oca:ocpp:uid:1:569236
	// Duration of the charging schedule in seconds. If the duration is left empty,
	// the last period will continue indefinitely or until end of the transaction if
	// chargingProfilePurpose = TxProfile.
	//
	Duration *int `json:"duration,omitempty" yaml:"duration,omitempty" mapstructure:"duration,omitempty"`

	// Identifies the ChargingSchedule.
	//
	Id int `json:"id" yaml:"id" mapstructure:"id"`

	// Charging_ Schedule. Min_ Charging_ Rate. Numeric
	// urn:x-oca:ocpp:uid:1:569239
	// Minimum charging rate supported by the EV. The unit of measure is defined by
	// the chargingRateUnit. This parameter is intended to be used by a local smart
	// charging algorithm to optimize the power allocation for in the case a charging
	// process is inefficient at lower charging rates. Accepts at most one digit
	// fraction (e.g. 8.1)
	//
	MinChargingRate *float64 `json:"minChargingRate,omitempty" yaml:"minChargingRate,omitempty" mapstructure:"minChargingRate,omitempty"`

	// SalesTariff corresponds to the JSON schema field "salesTariff".
	SalesTariff *SalesTariffType `json:"salesTariff,omitempty" yaml:"salesTariff,omitempty" mapstructure:"salesTariff,omitempty"`

	// Charging_ Schedule. Start_ Schedule. Date_ Time
	// urn:x-oca:ocpp:uid:1:569237
	// Starting point of an absolute schedule. If absent the schedule will be relative
	// to start of charging.
	//
	StartSchedule *string `json:"startSchedule,omitempty" yaml:"startSchedule,omitempty" mapstructure:"startSchedule,omitempty"`
}

type RecurrencyKindEnumType string

const RecurrencyKindEnumTypeDaily RecurrencyKindEnumType = "Daily"
const RecurrencyKindEnumTypeWeekly RecurrencyKindEnumType = "Weekly"

// Charging_ Profile
// urn:x-oca:ocpp:uid:2:233255
// A ChargingProfile consists of ChargingSchedule, describing the amount of power
// or current that can be delivered per time interval.
type ChargingProfileType struct {
	// ChargingProfileKind corresponds to the JSON schema field "chargingProfileKind".
	ChargingProfileKind ChargingProfileKindEnumType `json:"chargingProfileKind" yaml:"chargingProfileKind" mapstructure:"chargingProfileKind"`

	// ChargingProfilePurpose corresponds to the JSON schema field
	// "chargingProfilePurpose".
	ChargingProfilePurpose ChargingProfilePurposeEnumType `json:"chargingProfilePurpose" yaml:"chargingProfilePurpose" mapstructure:"chargingProfilePurpose"`

	// ChargingSchedule corresponds to the JSON schema field "chargingSchedule".
	ChargingSchedule []ChargingScheduleType `json:"chargingSchedule" yaml:"chargingSchedule" mapstructure:"chargingSchedule"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Identified_ Object. MRID. Numeric_ Identifier
	// urn:x-enexis:ecdm:uid:1:569198
	// Id of ChargingProfile.
	//
	Id int `json:"id" yaml:"id" mapstructure:"id"`

	// RecurrencyKind corresponds to the JSON schema field "recurrencyKind".
	RecurrencyKind *RecurrencyKindEnumType `json:"recurrencyKind,omitempty" yaml:"recurrencyKind,omitempty" mapstructure:"recurrencyKind,omitempty"`

	// Charging_ Profile. Stack_ Level. Counter
	// urn:x-oca:ocpp:uid:1:569230
	// Value determining level in hierarchy stack of profiles. Higher values have
	// precedence over lower values. Lowest level is 0.
	//
	StackLevel int `json:"stackLevel" yaml:"stackLevel" mapstructure:"stackLevel"`

	// SHALL only be included if ChargingProfilePurpose is set to TxProfile. The
	// transactionId is used to match the profile to a specific transaction.
	//
	TransactionId *string `json:"transactionId,omitempty" yaml:"transactionId,omitempty" mapstructure:"transactionId,omitempty"`

	// Charging_ Profile. Valid_ From. Date_ Time
	// urn:x-oca:ocpp:uid:1:569234
	// Point in time at which the profile starts to be valid. If absent, the profile
	// is valid as soon as it is received by the Charging Station.
	//
	ValidFrom *string `json:"validFrom,omitempty" yaml:"validFrom,omitempty" mapstructure:"validFrom,omitempty"`

	// Charging_ Profile. Valid_ To. Date_ Time
	// urn:x-oca:ocpp:uid:1:569235
	// Point in time at which the profile stops to be valid. If absent, the profile is
	// valid until it is replaced by another profile.
	//
	ValidTo *string `json:"validTo,omitempty" yaml:"validTo,omitempty" mapstructure:"validTo,omitempty"`
}

type RequestStartTransactionRequestJson struct {
	// ChargingProfile corresponds to the JSON schema field "chargingProfile".
	ChargingProfile *ChargingProfileType `json:"chargingProfile,omitempty" yaml:"chargingProfile,omitempty" mapstructure:"chargingProfile,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Number of the EVSE on which to start the transaction. EvseId SHALL be &gt; 0
	//
	EvseId *int `json:"evseId,omitempty" yaml:"evseId,omitempty" mapstructure:"evseId,omitempty"`

	// GroupIdToken corresponds to the JSON schema field "groupIdToken".
	GroupIdToken *IdTokenType `json:"groupIdToken,omitempty" yaml:"groupIdToken,omitempty" mapstructure:"groupIdToken,omitempty"`

	// IdToken corresponds to the JSON schema field "idToken".
	IdToken IdTokenType `json:"idToken" yaml:"idToken" mapstructure:"idToken"`

	// Id given by the server to this start request. The Charging Station might return
	// this in the &lt;&lt;transactioneventrequest, TransactionEventRequest&gt;&gt;,
	// letting the server know which transaction was started for this request. Use to
	// start a transaction.
	//
	RemoteStartId int `json:"remoteStartId" yaml:"remoteStartId" mapstructure:"remoteStartId"`
}

func (*RequestStartTransactionRequestJson) IsRequest() {}
