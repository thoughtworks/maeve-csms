package firestore

import (
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type token struct {
	CountryCode  string  `firestore:"country"`
	PartyId      string  `firestore:"partyId"`
	Type         string  `firestore:"type"`
	Uid          string  `firestore:"uid"`
	ContractId   string  `firestore:"contractId"`
	VisualNumber *string `firestore:"visual"`
	Issuer       string  `firestore:"issuer"`
	GroupId      *string `firestore:"group"`
	Valid        bool    `firestore:"valid"`
	LanguageCode *string `firestore:"lang"`
	CacheMode    string  `firestore:"cache"`
}

func (s *Store) SetToken(ctx context.Context, tok *store.Token) error {
	tokenRef := s.client.Doc(fmt.Sprintf("Token/%s", tok.Uid))
	tokenData := &token{
		CountryCode:  tok.CountryCode,
		PartyId:      tok.PartyId,
		Type:         tok.Type,
		Uid:          tok.Uid,
		ContractId:   tok.ContractId,
		VisualNumber: tok.VisualNumber,
		Issuer:       tok.Issuer,
		GroupId:      tok.GroupId,
		Valid:        tok.Valid,
		LanguageCode: tok.LanguageCode,
		CacheMode:    tok.CacheMode,
	}
	_, err := tokenRef.Set(ctx, tokenData)

	if err != nil {
		return fmt.Errorf("setting token: %s: %w", tok.Uid, err)
	}
	return nil
}

func (s *Store) LookupToken(ctx context.Context, tokenUid string) (*store.Token, error) {
	tokenRef := s.client.Doc(fmt.Sprintf("Token/%s", tokenUid))
	snap, err := tokenRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("lookup token %s: %w", tokenUid, err)
	}
	var tok token
	if err = snap.DataTo(&tok); err != nil {
		return nil, fmt.Errorf("map token: %s: %w", tokenUid, err)
	}
	return &store.Token{
		CountryCode:  tok.CountryCode,
		PartyId:      tok.PartyId,
		Type:         tok.Type,
		Uid:          tok.Uid,
		ContractId:   tok.ContractId,
		VisualNumber: tok.VisualNumber,
		Issuer:       tok.Issuer,
		GroupId:      tok.GroupId,
		Valid:        tok.Valid,
		LanguageCode: tok.LanguageCode,
		CacheMode:    tok.CacheMode,
	}, nil
}
