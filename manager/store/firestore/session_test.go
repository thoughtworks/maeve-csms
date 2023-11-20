// SPDX-License-Identifier: Apache-2.0

//go:build integration

package firestore_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"golang.org/x/net/context"
	"k8s.io/utils/clock"
	"testing"
)

func TestSetAndLookupSession(t *testing.T) {
	defer cleanupAllCollections(t, "myproject")

	ctx := context.Background()
	sessionStore, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	want := &store.Session{
		CountryCode:   "BEL",
		PartyId:       "TWK",
		Id:            "s001",
		StartDateTime: "", //Look at
		EndDateTime:   "",
		Kwh:           5,
		CdrToken: store.CdrToken{
			ContractId: "GBTWK012345678V",
			Type:       "RFID",
			Uid:        "MYRFIDTAG",
		},
		AuthMethod:  "AUTH_REQUEST", //may cause issue
		Currency:    "GBP",
		Status:      "ACTIVE",
		LastUpdated: "",
	}
	err = sessionStore.SetSession(ctx, want)
	require.NoError(t, err)

	got, err := sessionStore.LookupSession(ctx, "s001")
	require.NoError(t, err)

	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`, got.LastUpdated)
	got.LastUpdated = ""

	assert.Equal(t, want, got)
}

func TestListSessions(t *testing.T) {
	defer cleanupAllCollections(t, "myproject")

	ctx := context.Background()
	sessionStore, err := firestore.NewStore(ctx, "myproject", clock.RealClock{})
	require.NoError(t, err)

	sessions := make([]*store.Session, 20)
	for i := 0; i < 20; i++ {
		sessions[i] = &store.Session{
			CountryCode:   "BEL",
			PartyId:       "TWK",
			Id:            fmt.Sprintf("s%03d", i),
			StartDateTime: "", //Look at
			EndDateTime:   "",
			Kwh:           5,
			CdrToken: store.CdrToken{
				ContractId: "GBTWK012345678V",
				Type:       "RFID",
				Uid:        "MYRFIDTAG",
			},
			AuthMethod: "AUTH_REQUEST",
			Currency:   "GBP",
			Status:     "ACTIVE",
		}
	}

	for _, session := range sessions {
		err = sessionStore.SetSession(ctx, session)
		require.NoError(t, err)
	}

	got, err := sessionStore.ListSessions(ctx, 0, 10)
	require.NoError(t, err)

	assert.Equal(t, 10, len(got))
	for i, session := range got {
		session.LastUpdated = ""
		assert.Equal(t, sessions[i], got[i])
	}
}
