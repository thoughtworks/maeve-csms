// SPDX-License-Identifier: Apache-2.0

package registry

import "crypto/x509"

type SecurityProfile int

const (
	UnsecuredTransportWithBasicAuth SecurityProfile = iota
	TLSWithBasicAuth
	TLSWithClientSideCertificates
)

type ChargeStation struct {
	ClientId               string
	SecurityProfile        SecurityProfile
	Base64SHA256Password   string
	InvalidUsernameAllowed bool
}

type DeviceRegistry interface {
	LookupChargeStation(clientId string) (*ChargeStation, error)
	LookupCertificate(certHash string) (*x509.Certificate, error)
}
