// SPDX-License-Identifier: Apache-2.0

package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"k8s.io/utils/clock"
)

type Store struct {
	client *firestore.Client
	clock  clock.PassiveClock
}

func NewStore(ctx context.Context, gcloudProject string, clock clock.PassiveClock) (store.Engine, error) {
	client, err := firestore.NewClient(ctx, gcloudProject)
	if err != nil {
		return nil, fmt.Errorf("create new firestore client in %s: %w", gcloudProject, err)
	}

	return &Store{
		client: client,
		clock:  clock,
	}, nil
}
