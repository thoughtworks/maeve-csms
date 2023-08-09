// SPDX-License-Identifier: Apache-2.0

//go:build integration

package firestore_test

import (
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"testing"

	firestoreapi "cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
)

func deleteTokens(t *testing.T, gcloudProject string) {
	ctx := context.Background()

	client, err := firestoreapi.NewClient(ctx, gcloudProject)
	assert.NoError(t, err)

	col := client.Collection("Token")
	bulkwriter := client.BulkWriter(ctx)

	numDeleted := 0
	iter := col.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		assert.NoError(t, err)
		_, err = bulkwriter.Delete(doc.Ref)
		assert.NoError(t, err)
		numDeleted++
	}

	if numDeleted == 0 {
		bulkwriter.End()
	}

	bulkwriter.Flush()
}

func TestSetAndLookupToken(t *testing.T) {
	defer deleteTokens(t, "myproject")

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

	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`, got.LastUpdated)
	got.LastUpdated = ""

	assert.Equal(t, want, got)
}

func TestLookupTokenThatDoesNotExist(t *testing.T) {
	defer deleteTokens(t, "myproject")

	ctx := context.Background()

	tokenStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)

	got, err := tokenStore.LookupToken(ctx, "unknown-rfid")
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestListTokensWithNoMatches(t *testing.T) {
	defer deleteTokens(t, "myproject")

	ctx := context.Background()

	tokenStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)

	got, err := tokenStore.ListTokens(ctx, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, 0, len(got))
}

func TestListTokens(t *testing.T) {
	defer deleteTokens(t, "myproject")

	ctx := context.Background()

	tokenStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)

	contractId, err := ocpp.NormalizeEmaid("GB-TWK-C12345678")
	require.NoError(t, err)

	tokens := make([]*store.Token, 20)
	for i := 0; i < 20; i++ {
		tokens[i] = &store.Token{
			CountryCode: "GB",
			PartyId:     "TWK",
			Type:        "RFID",
			Uid:         fmt.Sprintf("123456%02d", i),
			ContractId:  contractId,
			Issuer:      "TWK",
			Valid:       true,
			CacheMode:   store.CacheModeAllowed,
		}
	}

	for _, token := range tokens {
		err = tokenStore.SetToken(ctx, token)
		require.NoError(t, err)
	}

	got, err := tokenStore.ListTokens(ctx, 0, 10)
	require.NoError(t, err)

	for _, token := range got {
		assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`, token.LastUpdated)
		token.LastUpdated = ""
	}

	require.Equal(t, 10, len(got))
	assert.Equal(t, tokens[:10], got)
}

func TestListTokensWithOffset(t *testing.T) {
	defer deleteTokens(t, "myproject")

	ctx := context.Background()

	tokenStore, err := firestore.NewStore(ctx, "myproject")
	require.NoError(t, err)

	contractId, err := ocpp.NormalizeEmaid("GB-TWK-C12345678")
	require.NoError(t, err)

	tokens := make([]*store.Token, 20)
	for i := 0; i < 20; i++ {
		tokens[i] = &store.Token{
			CountryCode: "GB",
			PartyId:     "TWK",
			Type:        "RFID",
			Uid:         fmt.Sprintf("123456%02d", i),
			ContractId:  contractId,
			Issuer:      "TWK",
			Valid:       true,
			CacheMode:   store.CacheModeAllowed,
		}
	}

	for _, token := range tokens {
		err = tokenStore.SetToken(ctx, token)
		require.NoError(t, err)
	}

	got, err := tokenStore.ListTokens(ctx, 5, 20)
	require.NoError(t, err)

	for _, token := range got {
		assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`, token.LastUpdated)
		token.LastUpdated = ""
	}

	require.Equal(t, 15, len(got))
	assert.Equal(t, tokens[5:20], got)
}
