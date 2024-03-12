// SPDX-License-Identifier: Apache-2.0

package ocpp201

type UnlockConnectorRequestJson struct {
	// This contains the identifier of the connector that needs to be unlocked.
	//
	ConnectorId int `json:"connectorId" yaml:"connectorId" mapstructure:"connectorId"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// This contains the identifier of the EVSE for which a connector needs to be
	// unlocked.
	//
	EvseId int `json:"evseId" yaml:"evseId" mapstructure:"evseId"`
}

func (*UnlockConnectorRequestJson) IsRequest() {}
