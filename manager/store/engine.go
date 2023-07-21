// SPDX-License-Identifier: Apache-2.0

package store

type Engine interface {
	ChargeStationAuthStore
	TokenStore
	TransactionStore
	CertificateStore
}
