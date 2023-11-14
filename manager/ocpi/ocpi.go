// SPDX-License-Identifier: Apache-2.0

package ocpi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	PushLocation(ctx context.Context, location Location) error
	PutSession(ctx context.Context, token store.Token) error
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
				Identifier: "commands",
				Role:       RECEIVER,
				Url:        fmt.Sprintf("%s/ocpi/receiver/2.2/commands", o.externalUrl),
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

func (o *OCPI) PushLocation(ctx context.Context, location Location) error {
	parties, err := o.store.ListPartyDetailsForRole(ctx, "EMSP")
	if err != nil {
		return err
	}
	for _, party := range parties {
		err = o.pushLocationToParty(ctx, party, location)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *OCPI) PushSession(ctx context.Context, session Session) error {
	parties, err := o.store.ListPartyDetailsForRole(ctx, "EMSP")
	if err != nil {
		return err
	}
	for _, party := range parties {
		err, _ = o.pushSessionToParty(ctx, party, session)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *OCPI) pushLocationToParty(ctx context.Context, party *store.OcpiParty, location Location) error {
	// TODO: retrieve endpoints from store, not via OCPI exchange
	versions, err := o.getVersions(ctx, party.Url, party.Token)
	if err != nil {
		return err
	}

	endpointUrl, err := getEndpointUrl(versions)
	if err != nil {
		return err
	}

	endpoints, err := o.getEndpoints(ctx, endpointUrl, party.Token)
	if err != nil {
		return err
	}

	locationsUrl, err := o.getLocationsUrl(endpoints)
	if err != nil {
		return err
	}

	err = o.putLocation(ctx, locationsUrl, party.CountryCode, party.PartyId, party.Token, location)
	if err != nil {
		return err
	}

	return nil
}

func (o *OCPI) pushSessionToParty(ctx context.Context, party *store.OcpiParty, session Session) (error, string) {
	// TODO: retrieve endpoints from store, not via OCPI exchange
	versions, err := o.getVersions(ctx, party.Url, party.Token)
	if err != nil {
		return err, ""
	}

	endpointUrl, err := getEndpointUrl(versions)
	if err != nil {
		return err, ""
	}

	endpoints, err := o.getEndpoints(ctx, endpointUrl, party.Token)
	if err != nil {
		return err, ""
	}

	sessionsUrl, err := o.getSessionsUrl(endpoints)
	if err != nil {
		return err, ""
	}
	//Call our function
	//err = o.putLocation(ctx, sessionsUrl, party.CountryCode, party.PartyId, party.Token, location)
	//if err != nil {
	//	return err
	//}

	return nil, sessionsUrl
}

func (o *OCPI) getLocationsUrl(endpoints []Endpoint) (string, error) {
	for _, endpoint := range endpoints {
		if endpoint.Identifier == "locations" && endpoint.Role == RECEIVER {
			return fmt.Sprintf("%s/%s/%s", endpoint.Url, o.countryCode, o.partyId), nil
		}
	}
	return "", errors.New("no locations endpoint for receiver found")
}

func (o *OCPI) putLocation(ctx context.Context, url string, toCountryCode string, toPartyId string, token string, location Location) error {
	b, err := json.Marshal(location)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/%s", url, location.Id), bytes.NewReader(b))
	if err != nil {
		return err
	}
	o.setRequestHeaders(ctx, req, token, toCountryCode, toPartyId)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

func (o *OCPI) setRequestHeaders(ctx context.Context, req *http.Request, token string, toCountryCode string, toPartyId string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	req.Header.Set("X-Request-ID", uuid.New().String())
	value, ok := ctx.Value(ContextKeyCorrelationId).(string)
	if !ok {
		value = uuid.New().String()
	}
	req.Header.Set("X-Correlation-ID", value)
	req.Header.Set("OCPI-from-country-code", o.countryCode)
	req.Header.Set("OCPI-from-party-id", o.partyId)
	req.Header.Set("OCPI-to-country-code", toCountryCode)
	req.Header.Set("OCPI-to-party-id", toPartyId)
}
func (o *OCPI) getSessionsUrl(endpoints []Endpoint) (string, error) {
	for _, endpoint := range endpoints {
		if endpoint.Identifier == "sessions" && endpoint.Role == RECEIVER {
			return fmt.Sprintf("%s/%s/%s", endpoint.Url, o.countryCode, o.partyId), nil
		}
	}
	return "", errors.New("no sessions endpoint for receiver found")
}

func (o *OCPI) PutSession(ctx context.Context, token store.Token) error {
	// will post the session to the emsp
	// we need partyId and CountryCode (how do we map chargeStationId with the emps that will need this info).
	// we'll use these to look up the party url and the api token
	// look at the token module to see if we can map the partyCode and CountryCode
	//generate a new sessionId
	session := new(Session)
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	url := ""
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/%s", url, session.Id), bytes.NewReader(b))
	if err != nil {
		return err
	}
	o.setRequestHeaders(ctx, req, token.Uid, token.CountryCode, token.PartyId)

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}
