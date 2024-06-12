package ocpp201

type SetChargingProfileRequestJson struct {
	EvseId          int
	ChargingProfile ChargingProfileType
}

func (*SetChargingProfileRequestJson) IsRequest() {}
