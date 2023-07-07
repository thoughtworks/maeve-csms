package bigtable_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twlabs/maeve-csms/manager/store"
	"github.com/twlabs/maeve-csms/manager/store/bigtable"
	"testing"
)

func TestSetAndLookupChargeStationAuth(t *testing.T) {
	ctx := context.Background()

	authStore, err := bigtable.NewStore(ctx, "myproject", "myinstance")
	require.NoError(t, err)
	assert.NotNil(t, authStore)

	want := &store.ChargeStationAuth{
		SecurityProfile:      store.TLSWithClientSideCertificates,
		Base64SHA256Password: "DEADBEEF",
	}
	err = authStore.SetChargeStationAuth(ctx, "cs001", want)
	require.NoError(t, err)

	got, err := authStore.LookupChargeStationAuth(ctx, "cs001")
	assert.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestLookupChargeStationAuthWithUnregisteredChargeStation(t *testing.T) {
	ctx := context.Background()

	authStore, err := bigtable.NewStore(ctx, "myproject", "myinstance")
	require.NoError(t, err)
	assert.NotNil(t, authStore)

	got, err := authStore.LookupChargeStationAuth(ctx, "cs002")
	assert.NoError(t, err)
	assert.Nil(t, got)
}
