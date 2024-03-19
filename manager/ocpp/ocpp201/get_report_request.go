// SPDX-License-Identifier: Apache-2.0

package ocpp201

type ComponentCriterionEnumType string

const ComponentCriterionEnumTypeActive ComponentCriterionEnumType = "Active"
const ComponentCriterionEnumTypeAvailable ComponentCriterionEnumType = "Available"
const ComponentCriterionEnumTypeEnabled ComponentCriterionEnumType = "Enabled"
const ComponentCriterionEnumTypeProblem ComponentCriterionEnumType = "Problem"

// Class to report components, variables and variable attributes and
// characteristics.
type ComponentVariableType struct {
	// Component corresponds to the JSON schema field "component".
	Component ComponentType `json:"component" yaml:"component" mapstructure:"component"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// Variable corresponds to the JSON schema field "variable".
	Variable *VariableType `json:"variable,omitempty" yaml:"variable,omitempty" mapstructure:"variable,omitempty"`
}

type GetReportRequestJson struct {
	// This field contains criteria for components for which a report is requested
	//
	ComponentCriteria []ComponentCriterionEnumType `json:"componentCriteria,omitempty" yaml:"componentCriteria,omitempty" mapstructure:"componentCriteria,omitempty"`

	// ComponentVariable corresponds to the JSON schema field "componentVariable".
	ComponentVariable []ComponentVariableType `json:"componentVariable,omitempty" yaml:"componentVariable,omitempty" mapstructure:"componentVariable,omitempty"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// The Id of the request.
	//
	RequestId int `json:"requestId" yaml:"requestId" mapstructure:"requestId"`
}

func (*GetReportRequestJson) IsRequest() {}
