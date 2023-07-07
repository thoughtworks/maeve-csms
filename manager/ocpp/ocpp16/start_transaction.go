// SPDX-License-Identifier: Apache-2.0

package ocpp16

type StartTransactionJson struct {
	// ConnectorId corresponds to the JSON schema field "connectorId".
	ConnectorId int `json:"connectorId" yaml:"connectorId" mapstructure:"connectorId"`

	// IdTag corresponds to the JSON schema field "idTag".
	IdTag string `json:"idTag" yaml:"idTag" mapstructure:"idTag"`

	// MeterStart corresponds to the JSON schema field "meterStart".
	MeterStart int `json:"meterStart" yaml:"meterStart" mapstructure:"meterStart"`

	// ReservationId corresponds to the JSON schema field "reservationId".
	ReservationId *int `json:"reservationId,omitempty" yaml:"reservationId,omitempty" mapstructure:"reservationId,omitempty"`

	// Timestamp corresponds to the JSON schema field "timestamp".
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
}

func (*StartTransactionJson) IsRequest() {}
