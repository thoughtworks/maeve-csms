// SPDX-License-Identifier: Apache-2.0

package ocpp201

type SetVariableResultType struct {
	// AttributeStatus corresponds to the JSON schema field "attributeStatus".
	AttributeStatus SetVariableStatusEnumType `json:"attributeStatus" yaml:"attributeStatus" mapstructure:"attributeStatus"`

	// AttributeStatusInfo corresponds to the JSON schema field "attributeStatusInfo".
	AttributeStatusInfo *StatusInfoType `json:"attributeStatusInfo,omitempty" yaml:"attributeStatusInfo,omitempty" mapstructure:"attributeStatusInfo,omitempty"`

	// AttributeType corresponds to the JSON schema field "attributeType".
	AttributeType *AttributeEnumType `json:"attributeType,omitempty" yaml:"attributeType,omitempty" mapstructure:"attributeType,omitempty"`

	// Component corresponds to the JSON schema field "component".
	Component ComponentType `json:"component" yaml:"component" mapstructure:"component"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Variable corresponds to the JSON schema field "variable".
	Variable VariableType `json:"variable" yaml:"variable" mapstructure:"variable"`
}

type SetVariableStatusEnumType string

const SetVariableStatusEnumTypeAccepted SetVariableStatusEnumType = "Accepted"
const SetVariableStatusEnumTypeNotSupportedAttributeType SetVariableStatusEnumType = "NotSupportedAttributeType"
const SetVariableStatusEnumTypeRebootRequired SetVariableStatusEnumType = "RebootRequired"
const SetVariableStatusEnumTypeRejected SetVariableStatusEnumType = "Rejected"
const SetVariableStatusEnumTypeUnknownComponent SetVariableStatusEnumType = "UnknownComponent"
const SetVariableStatusEnumTypeUnknownVariable SetVariableStatusEnumType = "UnknownVariable"

type SetVariablesResponseJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// SetVariableResult corresponds to the JSON schema field "setVariableResult".
	SetVariableResult []SetVariableResultType `json:"setVariableResult" yaml:"setVariableResult" mapstructure:"setVariableResult"`
}

func (*SetVariablesResponseJson) IsResponse() {}
