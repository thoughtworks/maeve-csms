package ocpp16

type StopTransactionJson struct {
	// IdTag corresponds to the JSON schema field "idTag".
	IdTag *string `json:"idTag,omitempty" yaml:"idTag,omitempty" mapstructure:"idTag,omitempty"`

	// MeterStop corresponds to the JSON schema field "meterStop".
	MeterStop int `json:"meterStop" yaml:"meterStop" mapstructure:"meterStop"`

	// Reason corresponds to the JSON schema field "reason".
	Reason *StopTransactionJsonReason `json:"reason,omitempty" yaml:"reason,omitempty" mapstructure:"reason,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`

	// TransactionData corresponds to the JSON schema field "transactionData".
	TransactionData []StopTransactionJsonTransactionDataElem `json:"transactionData,omitempty" yaml:"transactionData,omitempty" mapstructure:"transactionData,omitempty"`

	// TransactionId corresponds to the JSON schema field "transactionId".
	TransactionId int `json:"transactionId" yaml:"transactionId" mapstructure:"transactionId"`
}

func (*StopTransactionJson) IsRequest() {}

type StopTransactionJsonReason string

const StopTransactionJsonReasonDeAuthorized StopTransactionJsonReason = "DeAuthorized"
const StopTransactionJsonReasonEVDisconnected StopTransactionJsonReason = "EVDisconnected"
const StopTransactionJsonReasonEmergencyStop StopTransactionJsonReason = "EmergencyStop"
const StopTransactionJsonReasonHardReset StopTransactionJsonReason = "HardReset"
const StopTransactionJsonReasonLocal StopTransactionJsonReason = "Local"
const StopTransactionJsonReasonOther StopTransactionJsonReason = "Other"
const StopTransactionJsonReasonPowerLoss StopTransactionJsonReason = "PowerLoss"
const StopTransactionJsonReasonReboot StopTransactionJsonReason = "Reboot"
const StopTransactionJsonReasonRemote StopTransactionJsonReason = "Remote"
const StopTransactionJsonReasonSoftReset StopTransactionJsonReason = "SoftReset"
const StopTransactionJsonReasonUnlockCommand StopTransactionJsonReason = "UnlockCommand"

type StopTransactionJsonTransactionDataElem struct {
	// SampledValue corresponds to the JSON schema field "sampledValue".
	SampledValue []StopTransactionJsonTransactionDataElemSampledValueElem `json:"sampledValue" yaml:"sampledValue" mapstructure:"sampledValue"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
}

type StopTransactionJsonTransactionDataElemSampledValueElem struct {
	// Context corresponds to the JSON schema field "context".
	Context *StopTransactionJsonTransactionDataElemSampledValueElemContext `json:"context,omitempty" yaml:"context,omitempty" mapstructure:"context,omitempty"`

	// Format corresponds to the JSON schema field "format".
	Format *StopTransactionJsonTransactionDataElemSampledValueElemFormat `json:"format,omitempty" yaml:"format,omitempty" mapstructure:"format,omitempty"`

	// Location corresponds to the JSON schema field "location".
	Location *StopTransactionJsonTransactionDataElemSampledValueElemLocation `json:"location,omitempty" yaml:"location,omitempty" mapstructure:"location,omitempty"`

	// Measurand corresponds to the JSON schema field "measurand".
	Measurand *StopTransactionJsonTransactionDataElemSampledValueElemMeasurand `json:"measurand,omitempty" yaml:"measurand,omitempty" mapstructure:"measurand,omitempty"`

	// Phase corresponds to the JSON schema field "phase".
	Phase *StopTransactionJsonTransactionDataElemSampledValueElemPhase `json:"phase,omitempty" yaml:"phase,omitempty" mapstructure:"phase,omitempty"`

	// Unit corresponds to the JSON schema field "unit".
	Unit *StopTransactionJsonTransactionDataElemSampledValueElemUnit `json:"unit,omitempty" yaml:"unit,omitempty" mapstructure:"unit,omitempty"`

	// Value corresponds to the JSON schema field "value".
	Value string `json:"value" yaml:"value" mapstructure:"value"`
}

type StopTransactionJsonTransactionDataElemSampledValueElemContext string

const StopTransactionJsonTransactionDataElemSampledValueElemContextInterruptionBegin StopTransactionJsonTransactionDataElemSampledValueElemContext = "Interruption.Begin"
const StopTransactionJsonTransactionDataElemSampledValueElemContextInterruptionEnd StopTransactionJsonTransactionDataElemSampledValueElemContext = "Interruption.End"
const StopTransactionJsonTransactionDataElemSampledValueElemContextOther StopTransactionJsonTransactionDataElemSampledValueElemContext = "Other"
const StopTransactionJsonTransactionDataElemSampledValueElemContextSampleClock StopTransactionJsonTransactionDataElemSampledValueElemContext = "Sample.Clock"
const StopTransactionJsonTransactionDataElemSampledValueElemContextSamplePeriodic StopTransactionJsonTransactionDataElemSampledValueElemContext = "Sample.Periodic"
const StopTransactionJsonTransactionDataElemSampledValueElemContextTransactionBegin StopTransactionJsonTransactionDataElemSampledValueElemContext = "Transaction.Begin"
const StopTransactionJsonTransactionDataElemSampledValueElemContextTransactionEnd StopTransactionJsonTransactionDataElemSampledValueElemContext = "Transaction.End"
const StopTransactionJsonTransactionDataElemSampledValueElemContextTrigger StopTransactionJsonTransactionDataElemSampledValueElemContext = "Trigger"

type StopTransactionJsonTransactionDataElemSampledValueElemFormat string

const StopTransactionJsonTransactionDataElemSampledValueElemFormatRaw StopTransactionJsonTransactionDataElemSampledValueElemFormat = "Raw"
const StopTransactionJsonTransactionDataElemSampledValueElemFormatSignedData StopTransactionJsonTransactionDataElemSampledValueElemFormat = "SignedData"

type StopTransactionJsonTransactionDataElemSampledValueElemLocation string

const StopTransactionJsonTransactionDataElemSampledValueElemLocationBody StopTransactionJsonTransactionDataElemSampledValueElemLocation = "Body"
const StopTransactionJsonTransactionDataElemSampledValueElemLocationCable StopTransactionJsonTransactionDataElemSampledValueElemLocation = "Cable"
const StopTransactionJsonTransactionDataElemSampledValueElemLocationEV StopTransactionJsonTransactionDataElemSampledValueElemLocation = "EV"
const StopTransactionJsonTransactionDataElemSampledValueElemLocationInlet StopTransactionJsonTransactionDataElemSampledValueElemLocation = "Inlet"
const StopTransactionJsonTransactionDataElemSampledValueElemLocationOutlet StopTransactionJsonTransactionDataElemSampledValueElemLocation = "Outlet"

type StopTransactionJsonTransactionDataElemSampledValueElemMeasurand string

const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandCurrentExport StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Current.Export"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandCurrentImport StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Current.Import"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandCurrentOffered StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Current.Offered"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandEnergyActiveExportInterval StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Energy.Active.Export.Interval"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandEnergyActiveExportRegister StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Energy.Active.Export.Register"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandEnergyActiveImportInterval StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Energy.Active.Import.Interval"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandEnergyActiveImportRegister StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Energy.Active.Import.Register"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandEnergyReactiveExportInterval StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Energy.Reactive.Export.Interval"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandEnergyReactiveExportRegister StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Energy.Reactive.Export.Register"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandEnergyReactiveImportInterval StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Energy.Reactive.Import.Interval"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandEnergyReactiveImportRegister StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Energy.Reactive.Import.Register"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandFrequency StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Frequency"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandPowerActiveExport StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Power.Active.Export"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandPowerActiveImport StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Power.Active.Import"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandPowerFactor StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Power.Factor"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandPowerOffered StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Power.Offered"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandPowerReactiveExport StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Power.Reactive.Export"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandPowerReactiveImport StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Power.Reactive.Import"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandRPM StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "RPM"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandSoC StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "SoC"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandTemperature StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Temperature"
const StopTransactionJsonTransactionDataElemSampledValueElemMeasurandVoltage StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = "Voltage"

type StopTransactionJsonTransactionDataElemSampledValueElemPhase string

const StopTransactionJsonTransactionDataElemSampledValueElemPhaseL1 StopTransactionJsonTransactionDataElemSampledValueElemPhase = "L1"
const StopTransactionJsonTransactionDataElemSampledValueElemPhaseL1L2 StopTransactionJsonTransactionDataElemSampledValueElemPhase = "L1-L2"
const StopTransactionJsonTransactionDataElemSampledValueElemPhaseL1N StopTransactionJsonTransactionDataElemSampledValueElemPhase = "L1-N"
const StopTransactionJsonTransactionDataElemSampledValueElemPhaseL2 StopTransactionJsonTransactionDataElemSampledValueElemPhase = "L2"
const StopTransactionJsonTransactionDataElemSampledValueElemPhaseL2L3 StopTransactionJsonTransactionDataElemSampledValueElemPhase = "L2-L3"
const StopTransactionJsonTransactionDataElemSampledValueElemPhaseL2N StopTransactionJsonTransactionDataElemSampledValueElemPhase = "L2-N"
const StopTransactionJsonTransactionDataElemSampledValueElemPhaseL3 StopTransactionJsonTransactionDataElemSampledValueElemPhase = "L3"
const StopTransactionJsonTransactionDataElemSampledValueElemPhaseL3L1 StopTransactionJsonTransactionDataElemSampledValueElemPhase = "L3-L1"
const StopTransactionJsonTransactionDataElemSampledValueElemPhaseL3N StopTransactionJsonTransactionDataElemSampledValueElemPhase = "L3-N"
const StopTransactionJsonTransactionDataElemSampledValueElemPhaseN StopTransactionJsonTransactionDataElemSampledValueElemPhase = "N"

type StopTransactionJsonTransactionDataElemSampledValueElemUnit string

const StopTransactionJsonTransactionDataElemSampledValueElemUnitA StopTransactionJsonTransactionDataElemSampledValueElemUnit = "A"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitCelcius StopTransactionJsonTransactionDataElemSampledValueElemUnit = "Celcius"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitFahrenheit StopTransactionJsonTransactionDataElemSampledValueElemUnit = "Fahrenheit"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitK StopTransactionJsonTransactionDataElemSampledValueElemUnit = "K"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitKVA StopTransactionJsonTransactionDataElemSampledValueElemUnit = "kVA"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitKW StopTransactionJsonTransactionDataElemSampledValueElemUnit = "kW"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitKWh StopTransactionJsonTransactionDataElemSampledValueElemUnit = "kWh"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitKvar StopTransactionJsonTransactionDataElemSampledValueElemUnit = "kvar"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitKvarh StopTransactionJsonTransactionDataElemSampledValueElemUnit = "kvarh"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitPercent StopTransactionJsonTransactionDataElemSampledValueElemUnit = "Percent"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitV StopTransactionJsonTransactionDataElemSampledValueElemUnit = "V"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitVA StopTransactionJsonTransactionDataElemSampledValueElemUnit = "VA"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitVarh StopTransactionJsonTransactionDataElemSampledValueElemUnit = "varh"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitWh StopTransactionJsonTransactionDataElemSampledValueElemUnit = "Wh"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitW StopTransactionJsonTransactionDataElemSampledValueElemUnit = "W"
const StopTransactionJsonTransactionDataElemSampledValueElemUnitVar StopTransactionJsonTransactionDataElemSampledValueElemUnit = "var"

var enumValues_StopTransactionJsonTransactionDataElemSampledValueElemUnit = []interface{}{
	"Wh",
	"kWh",
	"varh",
	"kvarh",
	"W",
	"kW",
	"VA",
	"kVA",
	"var",
	"kvar",
	"A",
	"V",
	"K",
	"Celcius",
	"Fahrenheit",
	"Percent",
}

var enumValues_StopTransactionJsonTransactionDataElemSampledValueElemLocation = []interface{}{
	"Cable",
	"EV",
	"Inlet",
	"Outlet",
	"Body",
}

var enumValues_StopTransactionJsonTransactionDataElemSampledValueElemFormat = []interface{}{
	"Raw",
	"SignedData",
}

var enumValues_StopTransactionJsonTransactionDataElemSampledValueElemPhase = []interface{}{
	"L1",
	"L2",
	"L3",
	"N",
	"L1-N",
	"L2-N",
	"L3-N",
	"L1-L2",
	"L2-L3",
	"L3-L1",
}

var enumValues_StopTransactionJsonTransactionDataElemSampledValueElemMeasurand = []interface{}{
	"Energy.Active.Export.Register",
	"Energy.Active.Import.Register",
	"Energy.Reactive.Export.Register",
	"Energy.Reactive.Import.Register",
	"Energy.Active.Export.Interval",
	"Energy.Active.Import.Interval",
	"Energy.Reactive.Export.Interval",
	"Energy.Reactive.Import.Interval",
	"Power.Active.Export",
	"Power.Active.Import",
	"Power.Offered",
	"Power.Reactive.Export",
	"Power.Reactive.Import",
	"Power.Factor",
	"Current.Import",
	"Current.Export",
	"Current.Offered",
	"Voltage",
	"Frequency",
	"Temperature",
	"SoC",
	"RPM",
}

var enumValues_StopTransactionJsonTransactionDataElemSampledValueElemContext = []interface{}{
	"Interruption.Begin",
	"Interruption.End",
	"Sample.Clock",
	"Sample.Periodic",
	"Transaction.Begin",
	"Transaction.End",
	"Trigger",
	"Other",
}

var enumValues_StopTransactionJsonReason = []interface{}{
	"EmergencyStop",
	"EVDisconnected",
	"HardReset",
	"Local",
	"Other",
	"PowerLoss",
	"Reboot",
	"Remote",
	"SoftReset",
	"UnlockCommand",
	"DeAuthorized",
}
