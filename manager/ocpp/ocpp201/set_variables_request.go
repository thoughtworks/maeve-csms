// SPDX-License-Identifier: Apache-2.0

package ocpp201

type AttributeEnumType string

const AttributeEnumTypeActual AttributeEnumType = "Actual"
const AttributeEnumTypeTarget AttributeEnumType = "Target"
const AttributeEnumTypeMinSet AttributeEnumType = "MinSet"
const AttributeEnumTypeMaxSet AttributeEnumType = "MaxSet"

// A physical or logical component
type ComponentType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Evse corresponds to the JSON schema field "evse".
	Evse *EVSEType `json:"evse,omitempty" yaml:"evse,omitempty" mapstructure:"evse,omitempty"`

	// Name of instance in case the component exists as multiple instances. Case
	// Insensitive. strongly advised to use Camel Case.
	//
	Instance *string `json:"instance,omitempty" yaml:"instance,omitempty" mapstructure:"instance,omitempty"`

	// Name of the component. Name should be taken from the list of standardized
	// component names whenever possible. Case Insensitive. strongly advised to use
	// Camel Case.
	//
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

// Reference key to a component-variable.
type VariableType struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Name of instance in case the variable exists as multiple instances. Case
	// Insensitive. strongly advised to use Camel Case.
	//
	Instance *string `json:"instance,omitempty" yaml:"instance,omitempty" mapstructure:"instance,omitempty"`

	// Name of the variable. Name should be taken from the list of standardized
	// variable names whenever possible. Case Insensitive. strongly advised to use
	// Camel Case.
	//
	Name string `json:"name" yaml:"name" mapstructure:"name"`
}

type SetVariableDataType struct {
	// AttributeType corresponds to the JSON schema field "attributeType".
	AttributeType *AttributeEnumType `json:"attributeType,omitempty" yaml:"attributeType,omitempty" mapstructure:"attributeType,omitempty"`

	// Value to be assigned to attribute of variable.
	//
	// The Configuration Variable
	// &lt;&lt;configkey-configuration-value-size,ConfigurationValueSize&gt;&gt; can
	// be used to limit SetVariableData.attributeValue and
	// VariableCharacteristics.valueList. The max size of these values will always
	// remain equal.
	//
	AttributeValue string `json:"attributeValue" yaml:"attributeValue" mapstructure:"attributeValue"`

	// Component corresponds to the JSON schema field "component".
	Component ComponentType `json:"component" yaml:"component" mapstructure:"component"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Variable corresponds to the JSON schema field "variable".
	Variable VariableType `json:"variable" yaml:"variable" mapstructure:"variable"`
}

type SetVariablesRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// SetVariableData corresponds to the JSON schema field "setVariableData".
	SetVariableData []SetVariableDataType `json:"setVariableData" yaml:"setVariableData" mapstructure:"setVariableData"`
}

func (*SetVariablesRequestJson) IsRequest() {}
