package firestore_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"testing"
)

func TestSetAndLookupToken(t *testing.T) {
	ctx := context.Background()

	tokenStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)
	contractId, err := ocpp.NormalizeEmaid("GB-TWK-C12345678")
	require.NoError(t, err)
	want := &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "12345678",
		ContractId:  contractId,
		Issuer:      "TWK",
		Valid:       true,
		CacheMode:   store.CacheModeAllowed,
	}
	err = tokenStore.SetToken(ctx, want)
	require.NoError(t, err)

	got, err := tokenStore.LookupToken(ctx, "12345678")
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestLookupTokenThatDoesNotExist(t *testing.T) {
	ctx := context.Background()

	tokenStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)

	got, err := tokenStore.LookupToken(ctx, "unknown-rfid")
	require.NoError(t, err)
	require.Nil(t, got)
}
