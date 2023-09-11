// SPDX-License-Identifier: Apache-2.0

package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Store) SetLocation(ctx context.Context, loc *store.Location) error {
	locationRef := s.client.Doc(fmt.Sprintf("Location/%s", loc.Id))
	_, err := locationRef.Set(ctx, loc)
	if err != nil {
		return fmt.Errorf("setting location %s: %w", loc.Id, err)
	}
	return nil
}

func (s *Store) LookupLocation(ctx context.Context, locationId string) (*store.Location, error) {
	locationRef := s.client.Doc(fmt.Sprintf("Location/%s", locationId))
	snap, err := locationRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("lookup location %s: %w", locationId, err)
	}
	var location store.Location
	if err = snap.DataTo(&location); err != nil {
		return nil, fmt.Errorf("lookup location %s: %w", locationId, err)
	}
	location.LastUpdated = snap.UpdateTime.Format("2006-01-02T15:04:05Z")
	return &location, nil
}

func (s *Store) ListLocations(context context.Context, offset int, limit int) ([]*store.Location, error) {
	var locations []*store.Location
	iter := s.client.Collection("Location").OrderBy("Id", firestore.Asc).Offset(offset).Limit(limit).Documents(context)
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("next location: %w", err)
		}
		var loc store.Location
		if err = snap.DataTo(&loc); err != nil {
			return nil, fmt.Errorf("map location: %w", err)
		}
		loc.LastUpdated = snap.UpdateTime.Format("2006-01-02T15:04:05Z")
		locations = append(locations, &loc)
	}
	if locations == nil {
		locations = make([]*store.Location, 0)
	}
	return locations, nil
}
