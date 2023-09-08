// SPDX-License-Identifier: Apache-2.0

//go:build integration

package firestore_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/store/firestore"
	"k8s.io/utils/clock"
	"testing"
)

func TestNewStore(t *testing.T) {
	store, err := firestore.NewStore(context.Background(), "myproject", clock.RealClock{})
	require.NoError(t, err)
	assert.NotNil(t, store)
}
