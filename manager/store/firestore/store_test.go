package firestore_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"testing"
)

func TestNewStore(t *testing.T) {
	store, err := firestore.NewStore(context.Background(), "myproject")
	require.NoError(t, err)
	assert.NotNil(t, store)
}
