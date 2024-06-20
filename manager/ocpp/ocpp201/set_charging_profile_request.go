package ocpp201

type SetChargingProfileRequestJson struct {
	EvseId          int                 `json:"evseId"`
	ChargingProfile ChargingProfileType `json:"chargingProfile"`
}

func (*SetChargingProfileRequestJson) IsRequest() {}
