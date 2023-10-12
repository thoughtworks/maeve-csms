package firestore_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"golang.org/x/net/context"
	"k8s.io/utils/clock"
	"testing"
)

func TestSetAndLookupSession(t *testing.T) {
	//defer cleanupAllCollections(t, "myproject")

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
		LastUpdated: "2019-08-24T14:15:22Z",
	}
	err = sessionStore.SetSession(ctx, want)
	require.NoError(t, err)

	got, err := sessionStore.LookupSession(ctx, "s001")
	require.NoError(t, err)

	//assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)

	assert.Equal(t, want, got)
}
