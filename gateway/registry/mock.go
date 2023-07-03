package registry

type MockRegistry struct {
	ChargeStations map[string]*ChargeStation
}

func NewMockRegistry() *MockRegistry {
	return &MockRegistry{
		ChargeStations: make(map[string]*ChargeStation),
	}
}

func (m MockRegistry) LookupChargeStation(clientId string) *ChargeStation {
	return m.ChargeStations[clientId]
}
