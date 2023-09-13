// SPDX-License-Identifier: Apache-2.0

package firestore

import (
	"context"
	"errors"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"golang.org/x/exp/slog"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Store) SetRegistrationDetails(ctx context.Context, token string, registration *store.OcpiRegistration) error {
	slog.Info("setting registration", "token", token, "status", registration.Status)
	regRef := s.client.Doc(fmt.Sprintf("OcpiRegistration/%s", token))
	_, err := regRef.Set(ctx, registration)
	if err != nil {
		return fmt.Errorf("setting registration: %s: %w", token, err)
	}
	return nil
}

func (s *Store) GetRegistrationDetails(ctx context.Context, token string) (*store.OcpiRegistration, error) {
	regRef := s.client.Doc(fmt.Sprintf("OcpiRegistration/%s", token))
	snap, err := regRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
	}
	var registration store.OcpiRegistration
	err = snap.DataTo(&registration)
	if err != nil {
		return nil, fmt.Errorf("map registration %s: %w", token, err)
	}
	return &registration, nil
}

func (s *Store) DeleteRegistrationDetails(ctx context.Context, token string) error {
	regRef := s.client.Doc(fmt.Sprintf("OcpiRegistration/%s", token))
	_, err := regRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("delete registration %s: %w", token, err)
	}
	return nil
}

func (s *Store) SetPartyDetails(ctx context.Context, partyDetails *store.OcpiParty) error {
	partyRef := s.client.Doc(fmt.Sprintf("OcpiParty/%s/Id/%s:%s", partyDetails.Role, partyDetails.CountryCode, partyDetails.PartyId))
	_, err := partyRef.Set(ctx, partyDetails)
	if err != nil {
		return fmt.Errorf("setting party %s/%s:%s: %w", partyDetails.Role, partyDetails.CountryCode, partyDetails.PartyId, err)
	}
	return nil
}

func (s *Store) GetPartyDetails(ctx context.Context, role, countryCode, partyId string) (*store.OcpiParty, error) {
	partyRef := s.client.Doc(fmt.Sprintf("OcpiParty/%s/Id/%s:%s", role, countryCode, partyId))
	snap, err := partyRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("lookup party details %s/%s:%s: %w", role, countryCode, partyId, err)
	}
	var registration store.OcpiParty
	err = snap.DataTo(&registration)
	if err != nil {
		return nil, fmt.Errorf("map party details %s/%s:%s: %w", role, countryCode, partyId, err)
	}
	return &registration, nil
}

func (s *Store) ListPartyDetailsForRole(context context.Context, role string) ([]*store.OcpiParty, error) {
	var parties []*store.OcpiParty
	iter := s.client.Collection(fmt.Sprintf("OcpiParty/%s/Id", role)).Documents(context)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("next ocpiParty: %w", err)
		}
		var party store.OcpiParty
		if err = doc.DataTo(&party); err != nil {
			return nil, fmt.Errorf("map ocpiParty: %w", err)
		}
		parties = append(parties, &party)
	}
	if parties == nil {
		parties = make([]*store.OcpiParty, 0)
	}
	return parties, nil
}
