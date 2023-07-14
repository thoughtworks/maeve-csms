package store

type Engine interface {
	ChargeStationAuthStore
	TokenStore
	TransactionStore
}
