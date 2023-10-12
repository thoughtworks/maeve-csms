package store

import "context"

type Price struct {
	ExclVat float32
	InclVat float32
}

type CdrToken struct {
	ContractId string
	Type       string
	Uid        string
}

type Session struct {
	CountryCode   string
	PartyId       string
	Id            string
	StartDateTime string
	EndDateTime   string
	Kwh           float32
	CdrToken      CdrToken
	AuthMethod    string
	AuthReference *string
	LocationId    *string
	EvseId        *string
	ConnectorId   *string
	Currency      string
	TotalCost     Price //can be omitted/optional
	Status        string
	LastUpdated   string
}

type SessionStore interface {
	SetSession(ctx context.Context, session *Session) error
	LookupSession(ctx context.Context, sessionId string) (*Session, error)
	ListSessions(ctx context.Context, offset int, limit int) ([]*Session, error)
}
