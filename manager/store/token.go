package store

import "context"

const (
	CacheModeAlways         = "ALWAYS"
	CacheModeAllowed        = "ALLOWED"
	CacheModeAllowedOffline = "ALLOWED_OFFLINE"
	CacheModeNever          = "NEVER"
)

type Token struct {
	CountryCode  string
	PartyId      string
	Type         string
	Uid          string
	ContractId   string
	VisualNumber *string
	Issuer       string
	GroupId      *string
	Valid        bool
	LanguageCode *string
	CacheMode    string
	LastUpdated  string
}

type TokenStore interface {
	SetToken(ctx context.Context, token *Token) error
	LookupToken(ctx context.Context, tokenUid string) (*Token, error)
}
