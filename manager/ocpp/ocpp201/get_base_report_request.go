// SPDX-License-Identifier: Apache-2.0

package ocpp201

type ReportBaseEnumType string

const ReportBaseEnumTypeConfigurationInventory ReportBaseEnumType = "ConfigurationInventory"
const ReportBaseEnumTypeFullInventory ReportBaseEnumType = "FullInventory"
const ReportBaseEnumTypeSummaryInventory ReportBaseEnumType = "SummaryInventory"

type GetBaseReportRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// ReportBase corresponds to the JSON schema field "reportBase".
	ReportBase ReportBaseEnumType `json:"reportBase" yaml:"reportBase" mapstructure:"reportBase"`

	// The Id of the request.
	//
	RequestId int `json:"requestId" yaml:"requestId" mapstructure:"requestId"`
}

func (*GetBaseReportRequestJson) IsRequest() {}
