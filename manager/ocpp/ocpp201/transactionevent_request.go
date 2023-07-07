// SPDX-License-Identifier: Apache-2.0

package ocpp201

type ChargingStateEnumType string

const ChargingStateEnumTypeCharging ChargingStateEnumType = "Charging"
const ChargingStateEnumTypeEVConnected ChargingStateEnumType = "EVConnected"
const ChargingStateEnumTypeIdle ChargingStateEnumType = "Idle"
const ChargingStateEnumTypeSuspendedEV ChargingStateEnumType = "SuspendedEV"
const ChargingStateEnumTypeSuspendedEVSE ChargingStateEnumType = "SuspendedEVSE"

// EVSE
// urn:x-oca:ocpp:uid:2:233123
// Electric Vehicle Supply Equipment
type EVSEType struct {
	// An id to designate a specific connector (on an EVSE) by connector index number.
	//
	ConnectorId *int `json:"connectorId,omitempty" yaml:"connectorId,omitempty" mapstructure:"connectorId,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Identified_ Object. MRID. Numeric_ Identifier
	// urn:x-enexis:ecdm:uid:1:569198
	// EVSE Identifier. This contains a number (&gt; 0) designating an EVSE of the
	// Charging Station.
	//
	Id int `json:"id" yaml:"id" mapstructure:"id"`
}

type LocationEnumType string

const LocationEnumTypeBody LocationEnumType = "Body"
const LocationEnumTypeCable LocationEnumType = "Cable"
const LocationEnumTypeEV LocationEnumType = "EV"
const LocationEnumTypeInlet LocationEnumType = "Inlet"
const LocationEnumTypeOutlet LocationEnumType = "Outlet"

type MeasurandEnumType string

const MeasurandEnumTypeCurrentExport MeasurandEnumType = "Current.Export"
const MeasurandEnumTypeCurrentImport MeasurandEnumType = "Current.Import"
const MeasurandEnumTypeCurrentOffered MeasurandEnumType = "Current.Offered"
const MeasurandEnumTypeEnergyActiveExportInterval MeasurandEnumType = "Energy.Active.Export.Interval"
const MeasurandEnumTypeEnergyActiveExportRegister MeasurandEnumType = "Energy.Active.Export.Register"
const MeasurandEnumTypeEnergyActiveImportInterval MeasurandEnumType = "Energy.Active.Import.Interval"
const MeasurandEnumTypeEnergyActiveImportRegister MeasurandEnumType = "Energy.Active.Import.Register"
const MeasurandEnumTypeEnergyActiveNet MeasurandEnumType = "Energy.Active.Net"
const MeasurandEnumTypeEnergyApparentExport MeasurandEnumType = "Energy.Apparent.Export"
const MeasurandEnumTypeEnergyApparentImport MeasurandEnumType = "Energy.Apparent.Import"
const MeasurandEnumTypeEnergyApparentNet MeasurandEnumType = "Energy.Apparent.Net"
const MeasurandEnumTypeEnergyReactiveExportInterval MeasurandEnumType = "Energy.Reactive.Export.Interval"
const MeasurandEnumTypeEnergyReactiveExportRegister MeasurandEnumType = "Energy.Reactive.Export.Register"
const MeasurandEnumTypeEnergyReactiveImportInterval MeasurandEnumType = "Energy.Reactive.Import.Interval"
const MeasurandEnumTypeEnergyReactiveImportRegister MeasurandEnumType = "Energy.Reactive.Import.Register"
const MeasurandEnumTypeEnergyReactiveNet MeasurandEnumType = "Energy.Reactive.Net"
const MeasurandEnumTypeFrequency MeasurandEnumType = "Frequency"
const MeasurandEnumTypePowerActiveExport MeasurandEnumType = "Power.Active.Export"
const MeasurandEnumTypePowerActiveImport MeasurandEnumType = "Power.Active.Import"
const MeasurandEnumTypePowerFactor MeasurandEnumType = "Power.Factor"
const MeasurandEnumTypePowerOffered MeasurandEnumType = "Power.Offered"
const MeasurandEnumTypePowerReactiveExport MeasurandEnumType = "Power.Reactive.Export"
const MeasurandEnumTypePowerReactiveImport MeasurandEnumType = "Power.Reactive.Import"
const MeasurandEnumTypeSoC MeasurandEnumType = "SoC"
const MeasurandEnumTypeVoltage MeasurandEnumType = "Voltage"

// Meter_ Value
// urn:x-oca:ocpp:uid:2:233265
// Collection of one or more sampled values in MeterValuesRequest and
// TransactionEvent. All sampled values in a MeterValue are sampled at the same
// point in time.
type MeterValueType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// SampledValue corresponds to the JSON schema field "sampledValue".
	SampledValue []SampledValueType `json:"sampledValue" yaml:"sampledValue" mapstructure:"sampledValue"`

	// Meter_ Value. Timestamp. Date_ Time
	// urn:x-oca:ocpp:uid:1:569259
	// Timestamp for measured value(s).
	//
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
}

type PhaseEnumType string

const PhaseEnumTypeL1 PhaseEnumType = "L1"
const PhaseEnumTypeL1L2 PhaseEnumType = "L1-L2"
const PhaseEnumTypeL1N PhaseEnumType = "L1-N"
const PhaseEnumTypeL2 PhaseEnumType = "L2"
const PhaseEnumTypeL2L3 PhaseEnumType = "L2-L3"
const PhaseEnumTypeL2N PhaseEnumType = "L2-N"
const PhaseEnumTypeL3 PhaseEnumType = "L3"
const PhaseEnumTypeL3L1 PhaseEnumType = "L3-L1"
const PhaseEnumTypeL3N PhaseEnumType = "L3-N"
const PhaseEnumTypeN PhaseEnumType = "N"

type ReadingContextEnumType string

const ReadingContextEnumTypeInterruptionBegin ReadingContextEnumType = "Interruption.Begin"
const ReadingContextEnumTypeInterruptionEnd ReadingContextEnumType = "Interruption.End"
const ReadingContextEnumTypeOther ReadingContextEnumType = "Other"
const ReadingContextEnumTypeSampleClock ReadingContextEnumType = "Sample.Clock"
const ReadingContextEnumTypeSamplePeriodic ReadingContextEnumType = "Sample.Periodic"
const ReadingContextEnumTypeTransactionBegin ReadingContextEnumType = "Transaction.Begin"
const ReadingContextEnumTypeTransactionEnd ReadingContextEnumType = "Transaction.End"
const ReadingContextEnumTypeTrigger ReadingContextEnumType = "Trigger"

type ReasonEnumType string

const ReasonEnumTypeDeAuthorized ReasonEnumType = "DeAuthorized"
const ReasonEnumTypeEVDisconnected ReasonEnumType = "EVDisconnected"
const ReasonEnumTypeEmergencyStop ReasonEnumType = "EmergencyStop"
const ReasonEnumTypeEnergyLimitReached ReasonEnumType = "EnergyLimitReached"
const ReasonEnumTypeGroundFault ReasonEnumType = "GroundFault"
const ReasonEnumTypeImmediateReset ReasonEnumType = "ImmediateReset"
const ReasonEnumTypeLocal ReasonEnumType = "Local"
const ReasonEnumTypeLocalOutOfCredit ReasonEnumType = "LocalOutOfCredit"
const ReasonEnumTypeMasterPass ReasonEnumType = "MasterPass"
const ReasonEnumTypeOther ReasonEnumType = "Other"
const ReasonEnumTypeOvercurrentFault ReasonEnumType = "OvercurrentFault"
const ReasonEnumTypePowerLoss ReasonEnumType = "PowerLoss"
const ReasonEnumTypePowerQuality ReasonEnumType = "PowerQuality"
const ReasonEnumTypeReboot ReasonEnumType = "Reboot"
const ReasonEnumTypeRemote ReasonEnumType = "Remote"
const ReasonEnumTypeSOCLimitReached ReasonEnumType = "SOCLimitReached"
const ReasonEnumTypeStoppedByEV ReasonEnumType = "StoppedByEV"
const ReasonEnumTypeTimeLimitReached ReasonEnumType = "TimeLimitReached"
const ReasonEnumTypeTimeout ReasonEnumType = "Timeout"

// Sampled_ Value
// urn:x-oca:ocpp:uid:2:233266
// Single sampled value in MeterValues. Each value can be accompanied by optional
// fields.
//
// To save on mobile data usage, default values of all of the optional fields are
// such that. The value without any additional fields will be interpreted, as a
// register reading of active import energy in Wh (Watt-hour) units.
type SampledValueType struct {
	// Context corresponds to the JSON schema field "context".
	Context *ReadingContextEnumType `json:"context,omitempty" yaml:"context,omitempty" mapstructure:"context,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Location corresponds to the JSON schema field "location".
	Location *LocationEnumType `json:"location,omitempty" yaml:"location,omitempty" mapstructure:"location,omitempty"`

	// Measurand corresponds to the JSON schema field "measurand".
	Measurand *MeasurandEnumType `json:"measurand,omitempty" yaml:"measurand,omitempty" mapstructure:"measurand,omitempty"`

	// Phase corresponds to the JSON schema field "phase".
	Phase *PhaseEnumType `json:"phase,omitempty" yaml:"phase,omitempty" mapstructure:"phase,omitempty"`

	// SignedMeterValue corresponds to the JSON schema field "signedMeterValue".
	SignedMeterValue *SignedMeterValueType `json:"signedMeterValue,omitempty" yaml:"signedMeterValue,omitempty" mapstructure:"signedMeterValue,omitempty"`

	// UnitOfMeasure corresponds to the JSON schema field "unitOfMeasure".
	UnitOfMeasure *UnitOfMeasureType `json:"unitOfMeasure,omitempty" yaml:"unitOfMeasure,omitempty" mapstructure:"unitOfMeasure,omitempty"`

	// Sampled_ Value. Value. Measure
	// urn:x-oca:ocpp:uid:1:569260
	// Indicates the measured value.
	//
	//
	Value float64 `json:"value" yaml:"value" mapstructure:"value"`
}

// Represent a signed version of the meter value.
type SignedMeterValueType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Method used to encode the meter values before applying the digital signature
	// algorithm.
	//
	EncodingMethod string `json:"encodingMethod" yaml:"encodingMethod" mapstructure:"encodingMethod"`

	// Base64 encoded, sending depends on configuration variable
	// _PublicKeyWithSignedMeterValue_.
	//
	PublicKey string `json:"publicKey" yaml:"publicKey" mapstructure:"publicKey"`

	// Base64 encoded, contains the signed data which might contain more then just the
	// meter value. It can contain information like timestamps, reference to a
	// customer etc.
	//
	SignedMeterData string `json:"signedMeterData" yaml:"signedMeterData" mapstructure:"signedMeterData"`

	// Method used to create the digital signature.
	//
	SigningMethod string `json:"signingMethod" yaml:"signingMethod" mapstructure:"signingMethod"`
}

type TransactionEventEnumType string

const TransactionEventEnumTypeEnded TransactionEventEnumType = "Ended"
const TransactionEventEnumTypeStarted TransactionEventEnumType = "Started"
const TransactionEventEnumTypeUpdated TransactionEventEnumType = "Updated"

type TransactionEventRequestJson struct {
	// The maximum current of the connected cable in Ampere (A).
	//
	CableMaxCurrent *int `json:"cableMaxCurrent,omitempty" yaml:"cableMaxCurrent,omitempty" mapstructure:"cableMaxCurrent,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// EventType corresponds to the JSON schema field "eventType".
	EventType TransactionEventEnumType `json:"eventType" yaml:"eventType" mapstructure:"eventType"`

	// Evse corresponds to the JSON schema field "evse".
	Evse *EVSEType `json:"evse,omitempty" yaml:"evse,omitempty" mapstructure:"evse,omitempty"`

	// IdToken corresponds to the JSON schema field "idToken".
	IdToken *IdTokenType `json:"idToken,omitempty" yaml:"idToken,omitempty" mapstructure:"idToken,omitempty"`

	// MeterValue corresponds to the JSON schema field "meterValue".
	MeterValue []MeterValueType `json:"meterValue,omitempty" yaml:"meterValue,omitempty" mapstructure:"meterValue,omitempty"`

	// If the Charging Station is able to report the number of phases used, then it
	// SHALL provide it. When omitted the CSMS may be able to determine the number of
	// phases used via device management.
	//
	NumberOfPhasesUsed *int `json:"numberOfPhasesUsed,omitempty" yaml:"numberOfPhasesUsed,omitempty" mapstructure:"numberOfPhasesUsed,omitempty"`

	// Indication that this transaction event happened when the Charging Station was
	// offline. Default = false, meaning: the event occurred when the Charging Station
	// was online.
	//
	Offline bool `json:"offline,omitempty" yaml:"offline,omitempty" mapstructure:"offline,omitempty"`

	// This contains the Id of the reservation that terminates as a result of this
	// transaction.
	//
	ReservationId *int `json:"reservationId,omitempty" yaml:"reservationId,omitempty" mapstructure:"reservationId,omitempty"`

	// Incremental sequence number, helps with determining if all messages of a
	// transaction have been received.
	//
	SeqNo int `json:"seqNo" yaml:"seqNo" mapstructure:"seqNo"`

	// The date and time at which this transaction event occurred.
	//
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`

	// TransactionInfo corresponds to the JSON schema field "transactionInfo".
	TransactionInfo TransactionType `json:"transactionInfo" yaml:"transactionInfo" mapstructure:"transactionInfo"`

	// TriggerReason corresponds to the JSON schema field "triggerReason".
	TriggerReason TriggerReasonEnumType `json:"triggerReason" yaml:"triggerReason" mapstructure:"triggerReason"`
}

func (*TransactionEventRequestJson) IsRequest() {}

// Transaction
// urn:x-oca:ocpp:uid:2:233318
type TransactionType struct {
	// ChargingState corresponds to the JSON schema field "chargingState".
	ChargingState *ChargingStateEnumType `json:"chargingState,omitempty" yaml:"chargingState,omitempty" mapstructure:"chargingState,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// The ID given to remote start request (&lt;&lt;requeststarttransactionrequest,
	// RequestStartTransactionRequest&gt;&gt;. This enables to CSMS to match the
	// started transaction to the given start request.
	//
	RemoteStartId *int `json:"remoteStartId,omitempty" yaml:"remoteStartId,omitempty" mapstructure:"remoteStartId,omitempty"`

	// StoppedReason corresponds to the JSON schema field "stoppedReason".
	StoppedReason *ReasonEnumType `json:"stoppedReason,omitempty" yaml:"stoppedReason,omitempty" mapstructure:"stoppedReason,omitempty"`

	// Transaction. Time_ Spent_ Charging. Elapsed_ Time
	// urn:x-oca:ocpp:uid:1:569415
	// Contains the total time that energy flowed from EVSE to EV during the
	// transaction (in seconds). Note that timeSpentCharging is smaller or equal to
	// the duration of the transaction.
	//
	TimeSpentCharging *int `json:"timeSpentCharging,omitempty" yaml:"timeSpentCharging,omitempty" mapstructure:"timeSpentCharging,omitempty"`

	// This contains the Id of the transaction.
	//
	TransactionId string `json:"transactionId" yaml:"transactionId" mapstructure:"transactionId"`
}

type TriggerReasonEnumType string

const TriggerReasonEnumTypeAbnormalCondition TriggerReasonEnumType = "AbnormalCondition"
const TriggerReasonEnumTypeAuthorized TriggerReasonEnumType = "Authorized"
const TriggerReasonEnumTypeCablePluggedIn TriggerReasonEnumType = "CablePluggedIn"
const TriggerReasonEnumTypeChargingRateChanged TriggerReasonEnumType = "ChargingRateChanged"
const TriggerReasonEnumTypeChargingStateChanged TriggerReasonEnumType = "ChargingStateChanged"
const TriggerReasonEnumTypeDeauthorized TriggerReasonEnumType = "Deauthorized"
const TriggerReasonEnumTypeEVCommunicationLost TriggerReasonEnumType = "EVCommunicationLost"
const TriggerReasonEnumTypeEVConnectTimeout TriggerReasonEnumType = "EVConnectTimeout"
const TriggerReasonEnumTypeEVDeparted TriggerReasonEnumType = "EVDeparted"
const TriggerReasonEnumTypeEVDetected TriggerReasonEnumType = "EVDetected"
const TriggerReasonEnumTypeEnergyLimitReached TriggerReasonEnumType = "EnergyLimitReached"
const TriggerReasonEnumTypeMeterValueClock TriggerReasonEnumType = "MeterValueClock"
const TriggerReasonEnumTypeMeterValuePeriodic TriggerReasonEnumType = "MeterValuePeriodic"
const TriggerReasonEnumTypeRemoteStart TriggerReasonEnumType = "RemoteStart"
const TriggerReasonEnumTypeRemoteStop TriggerReasonEnumType = "RemoteStop"
const TriggerReasonEnumTypeResetCommand TriggerReasonEnumType = "ResetCommand"
const TriggerReasonEnumTypeSignedDataReceived TriggerReasonEnumType = "SignedDataReceived"
const TriggerReasonEnumTypeStopAuthorized TriggerReasonEnumType = "StopAuthorized"
const TriggerReasonEnumTypeTimeLimitReached TriggerReasonEnumType = "TimeLimitReached"
const TriggerReasonEnumTypeTrigger TriggerReasonEnumType = "Trigger"
const TriggerReasonEnumTypeUnlockCommand TriggerReasonEnumType = "UnlockCommand"

// Represents a UnitOfMeasure with a multiplier
type UnitOfMeasureType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Multiplier, this value represents the exponent to base 10. I.e. multiplier 3
	// means 10 raised to the 3rd power. Default is 0.
	//
	Multiplier int `json:"multiplier,omitempty" yaml:"multiplier,omitempty" mapstructure:"multiplier,omitempty"`

	// Unit of the value. Default = "Wh" if the (default) measurand is an "Energy"
	// type.
	// This field SHALL use a value from the list Standardized Units of Measurements
	// in Part 2 Appendices.
	// If an applicable unit is available in that list, otherwise a "custom" unit
	// might be used.
	//
	Unit string `json:"unit,omitempty" yaml:"unit,omitempty" mapstructure:"unit,omitempty"`
}
