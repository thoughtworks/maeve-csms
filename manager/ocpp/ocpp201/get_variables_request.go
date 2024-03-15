// SPDX-License-Identifier: Apache-2.0

package ocpp201

// Class to hold parameters for GetVariables request.
type GetVariableDataType struct {
	// AttributeType corresponds to the JSON schema field "attributeType".
	AttributeType *AttributeEnumType `json:"attributeType,omitempty" yaml:"attributeType,omitempty" mapstructure:"attributeType,omitempty"`

	// Component corresponds to the JSON schema field "component".
	Component ComponentType `json:"component" yaml:"component" mapstructure:"component"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Variable corresponds to the JSON schema field "variable".
	Variable VariableType `json:"variable" yaml:"variable" mapstructure:"variable"`
}

type GetVariablesRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// GetVariableData corresponds to the JSON schema field "getVariableData".
	GetVariableData []GetVariableDataType `json:"getVariableData" yaml:"getVariableData" mapstructure:"getVariableData"`
}

func (*GetVariablesRequestJson) IsRequest() {}
