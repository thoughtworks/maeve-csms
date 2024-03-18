// SPDX-License-Identifier: Apache-2.0

package ocpp201

type DataEnumType string

const DataEnumTypeBoolean DataEnumType = "boolean"
const DataEnumTypeDateTime DataEnumType = "dateTime"
const DataEnumTypeDecimal DataEnumType = "decimal"
const DataEnumTypeInteger DataEnumType = "integer"
const DataEnumTypeMemberList DataEnumType = "MemberList"
const DataEnumTypeOptionList DataEnumType = "OptionList"
const DataEnumTypeSequenceList DataEnumType = "SequenceList"
const DataEnumTypeString DataEnumType = "string"

type MutabilityEnumType string

const MutabilityEnumTypeReadOnly MutabilityEnumType = "ReadOnly"
const MutabilityEnumTypeReadWrite MutabilityEnumType = "ReadWrite"
const MutabilityEnumTypeWriteOnly MutabilityEnumType = "WriteOnly"

// Attribute data of a variable.
type VariableAttributeType struct {
	// If true, value that will never be changed by the Charging Station at runtime.
	// Default when omitted is false.
	//
	Constant bool `json:"constant,omitempty" yaml:"constant,omitempty" mapstructure:"constant,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Mutability corresponds to the JSON schema field "mutability".
	Mutability *MutabilityEnumType `json:"mutability,omitempty" yaml:"mutability,omitempty" mapstructure:"mutability,omitempty"`

	// If true, value will be persistent across system reboots or power down. Default
	// when omitted is false.
	//
	Persistent bool `json:"persistent,omitempty" yaml:"persistent,omitempty" mapstructure:"persistent,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type *AttributeEnumType `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`

	// Value of the attribute. May only be omitted when mutability is set to
	// 'WriteOnly'.
	//
	// The Configuration Variable
	// &lt;&lt;configkey-reporting-value-size,ReportingValueSize&gt;&gt; can be used
	// to limit GetVariableResult.attributeValue, VariableAttribute.value and
	// EventData.actualValue. The max size of these values will always remain equal.
	//
	Value *string `json:"value,omitempty" yaml:"value,omitempty" mapstructure:"value,omitempty"`
}

// Fixed read-only parameters of a variable.
type VariableCharacteristicsType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// DataType corresponds to the JSON schema field "dataType".
	DataType DataEnumType `json:"dataType" yaml:"dataType" mapstructure:"dataType"`

	// Maximum possible value of this variable. When the datatype of this Variable is
	// String, OptionList, SequenceList or MemberList, this field defines the maximum
	// length of the (CSV) string.
	//
	MaxLimit *float64 `json:"maxLimit,omitempty" yaml:"maxLimit,omitempty" mapstructure:"maxLimit,omitempty"`

	// Minimum possible value of this variable.
	//
	MinLimit *float64 `json:"minLimit,omitempty" yaml:"minLimit,omitempty" mapstructure:"minLimit,omitempty"`

	// Flag indicating if this variable supports monitoring.
	//
	SupportsMonitoring bool `json:"supportsMonitoring" yaml:"supportsMonitoring" mapstructure:"supportsMonitoring"`

	// Unit of the variable. When the transmitted value has a unit, this field SHALL
	// be included.
	//
	Unit *string `json:"unit,omitempty" yaml:"unit,omitempty" mapstructure:"unit,omitempty"`

	// Allowed values when variable is Option/Member/SequenceList.
	//
	// * OptionList: The (Actual) Variable value must be a single value from the
	// reported (CSV) enumeration list.
	//
	// * MemberList: The (Actual) Variable value  may be an (unordered) (sub-)set of
	// the reported (CSV) valid values list.
	//
	// * SequenceList: The (Actual) Variable value  may be an ordered (priority, etc)
	// (sub-)set of the reported (CSV) valid values.
	//
	// This is a comma separated list.
	//
	// The Configuration Variable
	// &lt;&lt;configkey-configuration-value-size,ConfigurationValueSize&gt;&gt; can
	// be used to limit SetVariableData.attributeValue and
	// VariableCharacteristics.valueList. The max size of these values will always
	// remain equal.
	//
	//
	ValuesList *string `json:"valuesList,omitempty" yaml:"valuesList,omitempty" mapstructure:"valuesList,omitempty"`
}

// Class to report components, variables and variable attributes and
// characteristics.
type ReportDataType struct {
	// Component corresponds to the JSON schema field "component".
	Component ComponentType `json:"component" yaml:"component" mapstructure:"component"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Variable corresponds to the JSON schema field "variable".
	Variable VariableType `json:"variable" yaml:"variable" mapstructure:"variable"`

	// VariableAttribute corresponds to the JSON schema field "variableAttribute".
	VariableAttribute []VariableAttributeType `json:"variableAttribute" yaml:"variableAttribute" mapstructure:"variableAttribute"`

	// VariableCharacteristics corresponds to the JSON schema field
	// "variableCharacteristics".
	VariableCharacteristics *VariableCharacteristicsType `json:"variableCharacteristics,omitempty" yaml:"variableCharacteristics,omitempty" mapstructure:"variableCharacteristics,omitempty"`
}

type NotifyReportRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Timestamp of the moment this message was generated at the Charging Station.
	//
	GeneratedAt string `json:"generatedAt" yaml:"generatedAt" mapstructure:"generatedAt"`

	// ReportData corresponds to the JSON schema field "reportData".
	ReportData []ReportDataType `json:"reportData,omitempty" yaml:"reportData,omitempty" mapstructure:"reportData,omitempty"`

	// The id of the GetReportRequest  or GetBaseReportRequest that requested this
	// report
	//
	RequestId int `json:"requestId" yaml:"requestId" mapstructure:"requestId"`

	// Sequence number of this message. First message starts at 0.
	//
	SeqNo int `json:"seqNo" yaml:"seqNo" mapstructure:"seqNo"`

	// “to be continued” indicator. Indicates whether another part of the report
	// follows in an upcoming notifyReportRequest message. Default value when omitted
	// is false.
	//
	//
	Tbc bool `json:"tbc,omitempty" yaml:"tbc,omitempty" mapstructure:"tbc,omitempty"`
}

func (*NotifyReportRequestJson) IsRequest() {}
