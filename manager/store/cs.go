package store

import "context"

type SecurityProfile int8

const (
	UnsecuredTransportWithBasicAuth SecurityProfile = iota
	TLSWithBasicAuth
	TLSWithClientSideCertificates
)

type ChargeStationAuth struct {
	SecurityProfile      SecurityProfile
	Base64SHA256Password string
}

type ChargeStationAuthStore interface {
	SetChargeStationAuth(ctx context.Context, chargeStationId string, auth *ChargeStationAuth) error
	LookupChargeStationAuth(ctx context.Context, chargeStationId string) (*ChargeStationAuth, error)
}
