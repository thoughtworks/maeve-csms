// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
)

type OcpiRegistrationStatusType string

var (
	OcpiRegistrationStatusRegistered OcpiRegistrationStatusType = "registered"
	OcpiRegistrationStatusPending    OcpiRegistrationStatusType = "pending"
)

type OcpiRegistration struct {
	Status OcpiRegistrationStatusType
}

type OcpiParty struct {
	CountryCode string
	PartyId     string
	Role        string
	Url         string
	Token       string
}

type OcpiStore interface {
	SetRegistrationDetails(ctx context.Context, token string, registration *OcpiRegistration) error
	GetRegistrationDetails(ctx context.Context, token string) (*OcpiRegistration, error)
	DeleteRegistrationDetails(ctx context.Context, token string) error

	SetPartyDetails(ctx context.Context, partyDetails *OcpiParty) error
	GetPartyDetails(ctx context.Context, role, countryCode, partyId string) (*OcpiParty, error)
	ListPartyDetailsForRole(ctx context.Context, role string) ([]*OcpiParty, error)
}
