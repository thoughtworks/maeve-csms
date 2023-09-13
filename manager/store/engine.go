// SPDX-License-Identifier: Apache-2.0

package store

type Engine interface {
	ChargeStationAuthStore
	ChargeStationSettingsStore
	ChargeStationRuntimeDetailsStore
	ChargeStationInstallCertificatesStore
	TokenStore
	TransactionStore
	CertificateStore
	OcpiStore
	LocationStore
}
