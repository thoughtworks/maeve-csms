// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"time"
)

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

type ChargeStationSettingStatus string

var (
	ChargeStationSettingStatusPending        ChargeStationSettingStatus = "Pending"
	ChargeStationSettingStatusAccepted       ChargeStationSettingStatus = "Accepted"
	ChargeStationSettingStatusRejected       ChargeStationSettingStatus = "Rejected"
	ChargeStationSettingStatusRebootRequired ChargeStationSettingStatus = "RebootRequired"
	ChargeStationSettingStatusNotSupported   ChargeStationSettingStatus = "NotSupported"
)

type ChargeStationSetting struct {
	Value       string
	Status      ChargeStationSettingStatus
	LastUpdated time.Time
}

type ChargeStationSettings struct {
	ChargeStationId string
	Settings        map[string]*ChargeStationSetting
}

type ChargeStationSettingsStore interface {
	UpdateChargeStationSettings(ctx context.Context, chargeStationId string, settings *ChargeStationSettings) error
	LookupChargeStationSettings(ctx context.Context, chargeStationId string) (*ChargeStationSettings, error)
	ListChargeStationSettings(ctx context.Context, pageSize int, previousChargeStationId string) ([]*ChargeStationSettings, error)
}

type ChargeStationRuntimeDetails struct {
	OcppVersion string
}

type ChargeStationRuntimeDetailsStore interface {
	SetChargeStationRuntimeDetails(ctx context.Context, chargeStationId string, details *ChargeStationRuntimeDetails) error
	LookupChargeStationRuntimeDetails(ctx context.Context, chargeStationId string) (*ChargeStationRuntimeDetails, error)
}
