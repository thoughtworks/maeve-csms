package ocpp201

type ConnectorStatusEnumType string

const ConnectorStatusEnumTypeAvailable ConnectorStatusEnumType = "Available"
const ConnectorStatusEnumTypeFaulted ConnectorStatusEnumType = "Faulted"
const ConnectorStatusEnumTypeOccupied ConnectorStatusEnumType = "Occupied"
const ConnectorStatusEnumTypeReserved ConnectorStatusEnumType = "Reserved"
const ConnectorStatusEnumTypeUnavailable ConnectorStatusEnumType = "Unavailable"

type StatusNotificationRequestJson struct {
	// The id of the connector within the EVSE for which the status is reported.
	//
	ConnectorId int `json:"connectorId" yaml:"connectorId" mapstructure:"connectorId"`

	// ConnectorStatus corresponds to the JSON schema field "connectorStatus".
	ConnectorStatus ConnectorStatusEnumType `json:"connectorStatus" yaml:"connectorStatus" mapstructure:"connectorStatus"`

	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// The id of the EVSE to which the connector belongs for which the the status is
	// reported.
	//
	EvseId int `json:"evseId" yaml:"evseId" mapstructure:"evseId"`

	// The time for which the status is reported. If absent time of receipt of the
	// message will be assumed.
	//
	Timestamp string `json:"timestamp" yaml:"timestamp" mapstructure:"timestamp"`
}

func (*StatusNotificationRequestJson) IsRequest() {}
