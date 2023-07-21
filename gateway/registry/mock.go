// SPDX-License-Identifier: Apache-2.0

package registry

import "crypto/x509"

type MockRegistry struct {
	ChargeStations map[string]*ChargeStation
	Certificates   map[string]*x509.Certificate
}

func NewMockRegistry() *MockRegistry {
	return &MockRegistry{
		ChargeStations: make(map[string]*ChargeStation),
		Certificates:   make(map[string]*x509.Certificate),
	}
}

func (m MockRegistry) LookupChargeStation(clientId string) (*ChargeStation, error) {
	return m.ChargeStations[clientId], nil
}

func (m MockRegistry) LookupCertificate(certHash string) (*x509.Certificate, error) {
	return m.Certificates[certHash], nil
}
