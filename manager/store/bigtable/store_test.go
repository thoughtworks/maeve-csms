package bigtable_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twlabs/maeve-csms/manager/store/bigtable"
	"testing"
)

func TestNewStore(t *testing.T) {
	ctx := context.Background()

	authStore, err := bigtable.NewStore(ctx, "myproject", "myinstance")
	require.NoError(t, err)
	assert.NotNil(t, authStore)
}
