// SPDX-License-Identifier: Apache-2.0

package ocpi

import (
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"net/http"
)

//go:generate oapi-codegen -config cfg.yaml ocpi22-spec.yaml

type Api interface {
	SetExternalUrl(serverUrl string)
	RegisterNewParty(ctx context.Context, url, token string) error
	GetVersions(ctx context.Context) ([]Version, error)
	GetVersion(ctx context.Context) (VersionDetail, error)
	SetCredentials(ctx context.Context, token string, credentials Credentials) error
	SetToken(ctx context.Context, token Token) error
	GetToken(ctx context.Context, countryCode string, partyID string, tokenUID string) (*Token, error)
}

type OCPI struct {
	store       store.Engine
	httpClient  *http.Client
	externalUrl string
	countryCode string
	partyId     string
}

func NewOCPI(store store.Engine, httpClient *http.Client, countryCode, partyId string) *OCPI {
	return &OCPI{
		store:       store,
		httpClient:  httpClient,
		countryCode: countryCode,
		partyId:     partyId,
	}
}

func (o *OCPI) SetExternalUrl(externalUrl string) {
	o.externalUrl = externalUrl
}

func (o *OCPI) GetVersions(context.Context) ([]Version, error) {
	return []Version{
		{
			Url:     fmt.Sprintf("%s/ocpi/2.2", o.externalUrl),
			Version: "2.2",
		},
	}, nil
}

func (o *OCPI) GetVersion(context.Context) (VersionDetail, error) {
	return VersionDetail{
		Endpoints: []Endpoint{
			{
				Identifier: "credentials",
				Role:       RECEIVER,
				Url:        fmt.Sprintf("%s/ocpi/2.2/credentials", o.externalUrl),
			},
			{
				Identifier: "tokens",
				Role:       RECEIVER,
				Url:        fmt.Sprintf("%s/ocpi/receiver/2.2/tokens/", o.externalUrl),
			},
		},
		Version: "2.2",
	}, nil
}

func (o *OCPI) SetCredentials(ctx context.Context, token string, credentials Credentials) error {
	for _, role := range credentials.Roles {
		err := o.store.SetPartyDetails(ctx, &store.OcpiParty{
			Role:        string(role.Role),
			CountryCode: role.CountryCode,
			PartyId:     role.PartyId,
			Url:         credentials.Url,
			Token:       credentials.Token,
		})
		if err != nil {
			return err
		}
	}

	reg, err := o.store.GetRegistrationDetails(ctx, token)
	if err != nil {
		return err
	}

	if reg != nil && reg.Status == store.OcpiRegistrationStatusPending {
		// register new party
		err = o.RegisterNewParty(ctx, credentials.Url, credentials.Token)
		if err != nil {
			return err
		}
		// delete old token
		err := o.store.DeleteRegistrationDetails(ctx, token)
		if err != nil {
			return err
		}
	}

	// store new token
	err = o.store.SetRegistrationDetails(ctx, credentials.Token, &store.OcpiRegistration{
		Status: store.OcpiRegistrationStatusRegistered,
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *OCPI) GetToken(ctx context.Context, countryCode string, partyID string, tokenUID string) (*Token, error) {
	tok, err := o.store.LookupToken(ctx, tokenUID)
	if err != nil {
		return nil, err
	}
	if tok.CountryCode != countryCode || tok.PartyId != partyID {
		return nil, nil
	}
	return &Token{
		ContractId:   tok.ContractId,
		CountryCode:  tok.CountryCode,
		GroupId:      tok.GroupId,
		Issuer:       tok.Issuer,
		Language:     tok.LanguageCode,
		LastUpdated:  tok.LastUpdated,
		PartyId:      tok.PartyId,
		Type:         TokenType(tok.Type),
		Uid:          tok.Uid,
		Valid:        tok.Valid,
		VisualNumber: tok.VisualNumber,
		Whitelist:    TokenWhitelist(tok.CacheMode),
	}, nil
}

func (o *OCPI) SetToken(ctx context.Context, token Token) error {
	tok := &store.Token{
		CountryCode:  token.CountryCode,
		PartyId:      token.PartyId,
		Type:         string(token.Type),
		Uid:          token.Uid,
		ContractId:   token.ContractId,
		VisualNumber: token.VisualNumber,
		Issuer:       token.Issuer,
		GroupId:      token.GroupId,
		Valid:        token.Valid,
		LanguageCode: token.Language,
		CacheMode:    string(token.Whitelist),
	}

	return o.store.SetToken(ctx, tok)
}
