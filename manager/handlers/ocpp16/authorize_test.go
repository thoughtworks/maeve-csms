// SPDX-License-Identifier: Apache-2.0

package ocpp16_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"testing"
	"time"
)

func TestAuthorizeKnownRfidCard(t *testing.T) {
	engine := inmemory.NewStore()
	err := engine.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "MYRFIDCARD",
		ContractId:  "GBTWK012345678V",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "NEVER",
		LastUpdated: time.Now().Format(time.RFC3339),
	})
	require.NoError(t, err)

	ah := handlers.AuthorizeHandler{
		TokenStore: engine,
	}

	req := &types.AuthorizeJson{
		IdTag: "MYRFIDCARD",
	}

	got, err := ah.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.AuthorizeResponseJson{
		IdTagInfo: types.AuthorizeResponseJsonIdTagInfo{
			Status: types.AuthorizeResponseJsonIdTagInfoStatusAccepted,
		},
	}

	assert.Equal(t, want, got)
}

func TestAuthorizeWithUnknownRfidCard(t *testing.T) {
	engine := inmemory.NewStore()

	ah := handlers.AuthorizeHandler{
		TokenStore: engine,
	}

	req := &types.AuthorizeJson{
		IdTag: "MYBADRFID",
	}

	got, err := ah.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.AuthorizeResponseJson{
		IdTagInfo: types.AuthorizeResponseJsonIdTagInfo{
			Status: types.AuthorizeResponseJsonIdTagInfoStatusInvalid,
		},
	}

	assert.Equal(t, want, got)
}
