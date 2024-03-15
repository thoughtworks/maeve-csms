// SPDX-License-Identifier: Apache-2.0

package ocpp201

// Class to hold results of GetVariables request.
type GetVariableResultType struct {
	// AttributeStatus corresponds to the JSON schema field "attributeStatus".
	AttributeStatus GetVariableStatusEnumType `json:"attributeStatus" yaml:"attributeStatus" mapstructure:"attributeStatus"`

	// AttributeStatusInfo corresponds to the JSON schema field "attributeStatusInfo".
	AttributeStatusInfo *StatusInfoType `json:"attributeStatusInfo,omitempty" yaml:"attributeStatusInfo,omitempty" mapstructure:"attributeStatusInfo,omitempty"`

	// AttributeType corresponds to the JSON schema field "attributeType".
	AttributeType *AttributeEnumType `json:"attributeType,omitempty" yaml:"attributeType,omitempty" mapstructure:"attributeType,omitempty"`

	// Value of requested attribute type of component-variable. This field can only be
	// empty when the given status is NOT accepted.
	//
	// The Configuration Variable
	// &lt;&lt;configkey-reporting-value-size,ReportingValueSize&gt;&gt; can be used
	// to limit GetVariableResult.attributeValue, VariableAttribute.value and
	// EventData.actualValue. The max size of these values will always remain equal.
	//
	//
	AttributeValue *string `json:"attributeValue,omitempty" yaml:"attributeValue,omitempty" mapstructure:"attributeValue,omitempty"`

	// Component corresponds to the JSON schema field "component".
	Component ComponentType `json:"component" yaml:"component" mapstructure:"component"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Variable corresponds to the JSON schema field "variable".
	Variable VariableType `json:"variable" yaml:"variable" mapstructure:"variable"`
}

type GetVariableStatusEnumType string

const GetVariableStatusEnumTypeAccepted GetVariableStatusEnumType = "Accepted"
const GetVariableStatusEnumTypeNotSupportedAttributeType GetVariableStatusEnumType = "NotSupportedAttributeType"
const GetVariableStatusEnumTypeRejected GetVariableStatusEnumType = "Rejected"
const GetVariableStatusEnumTypeUnknownComponent GetVariableStatusEnumType = "UnknownComponent"
const GetVariableStatusEnumTypeUnknownVariable GetVariableStatusEnumType = "UnknownVariable"

type GetVariablesResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// GetVariableResult corresponds to the JSON schema field "getVariableResult".
	GetVariableResult []GetVariableResultType `json:"getVariableResult" yaml:"getVariableResult" mapstructure:"getVariableResult"`
}

func (*GetVariablesResponseJson) IsResponse() {}
