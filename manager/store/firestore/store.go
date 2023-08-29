// SPDX-License-Identifier: Apache-2.0

package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
)

type Store struct {
	client *firestore.Client
}

func NewStore(ctx context.Context, gcloudProject string) (store.Engine, error) {
	client, err := firestore.NewClient(ctx, gcloudProject)
	if err != nil {
		return nil, fmt.Errorf("create new firestore client in %s: %w", gcloudProject, err)
	}

	return &Store{
		client: client,
	}, nil
}
