package ocpp201

type SetChargingProfileRequestJson struct {
	EvseId              int
	ChargingProfileType ChargingProfileType
}

func (*SetChargingProfileRequestJson) IsRequest() {}
