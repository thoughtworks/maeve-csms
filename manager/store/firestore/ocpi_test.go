// SPDX-License-Identifier: Apache-2.0

//go:build integration

package firestore_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"k8s.io/utils/clock"
	"testing"
)

func TestSetAndLookupRegistrationDetails(t *testing.T) {
	ctx := context.Background()

	engine, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	token := "abcdef123456"
	want := &store.OcpiRegistration{
		Status: store.OcpiRegistrationStatusRegistered,
	}

	err = engine.SetRegistrationDetails(ctx, token, want)
	require.NoError(t, err)

	got, err := engine.GetRegistrationDetails(ctx, token)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestDeleteRegistrationDetails(t *testing.T) {
	ctx := context.Background()

	engine, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	token := "abcdef123456"
	stored := &store.OcpiRegistration{
		Status: store.OcpiRegistrationStatusRegistered,
	}

	err = engine.SetRegistrationDetails(ctx, token, stored)
	require.NoError(t, err)

	err = engine.DeleteRegistrationDetails(ctx, token)
	require.NoError(t, err)

	got, err := engine.GetRegistrationDetails(ctx, token)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestSetAndLookupPartyDetails(t *testing.T) {
	ctx := context.Background()

	engine, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	want := &store.OcpiParty{
		Role:        "CPO",
		CountryCode: "GB",
		PartyId:     "TWK",
		Url:         "https://example.com/ocpi/versions",
		Token:       "abcdef123456",
	}

	err = engine.SetPartyDetails(ctx, want)
	require.NoError(t, err)

	got, err := engine.GetPartyDetails(ctx, want.Role, want.CountryCode, want.PartyId)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestSetAndListPartyDetails(t *testing.T) {
	ctx := context.Background()

	engine, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	want := &store.OcpiParty{
		Role:        "EMSP",
		CountryCode: "GB",
		PartyId:     "TWK",
		Url:         "https://example.com/ocpi/versions",
		Token:       "abcdef123456",
	}

	err = engine.SetPartyDetails(ctx, want)
	require.NoError(t, err)

	got, err := engine.ListPartyDetailsForRole(ctx, "EMSP")
	require.NoError(t, err)

	assert.Equal(t, 1, len(got))
	assert.Equal(t, want, got[0])
}
