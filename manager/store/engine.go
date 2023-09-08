// SPDX-License-Identifier: Apache-2.0

package store

type Engine interface {
	ChargeStationAuthStore
	ChargeStationSettingsStore
	ChargeStationRuntimeDetailsStore
	TokenStore
	TransactionStore
	CertificateStore
	OcpiStore
}
