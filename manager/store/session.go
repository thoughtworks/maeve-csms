package store

import "context"

type Price struct {
	ExclVat float32
	InclVat float32
}

type CdrToken struct {
	ContractId string
	Type       CdrTokenType
	Uid        string
}

type ChargingPeriod struct {
	Dimensions    []CdrDimension `json:"dimensions"`
	StartDateTime string         `json:"start_date_time"`
	TariffId      *string        `json:"tariff_id,omitempty"`
}

type CdrDimension struct {
	Type   CdrDimensionType `json:"type"`
	Volume float32          `json:"volume"`
}
type CdrTokenType string
type CdrDimensionType string
type SessionStatus string

type SessionAuthMethod string

type Session struct {
	AuthMethod             SessionAuthMethod
	AuthorizationReference *string
	CdrToken               CdrToken
	ChargingPeriods        *[]ChargingPeriod
	ConnectorId            string
	CountryCode            string
	Currency               string
	EndDateTime            *string
	EvseUid                string
	Id                     string
	Kwh                    float32
	LastUpdated            string
	LocationId             string
	MeterId                *string
	PartyId                string
	StartDateTime          string
	Status                 SessionStatus
	TotalCost              *Price
}

type SessionStore interface {
	SetSession(ctx context.Context, session *Session) error
	LookupSession(ctx context.Context, sessionId string) (*Session, error)
	ListSessions(ctx context.Context, offset int, limit int) ([]*Session, error)
}
