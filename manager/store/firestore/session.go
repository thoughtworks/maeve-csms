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

func (s *Store) SetSession(ctx context.Context, session *store.Session) error {
	sessionRef := s.client.Doc(fmt.Sprintf("Session/%s", session.Id))
	_, err := sessionRef.Set(ctx, session)
	if err != nil {
		return fmt.Errorf("setting session %s: %w", session.Id, err)
	}
	return nil
}

func (s *Store) LookupSession(ctx context.Context, sessionId string) (*store.Session, error) {
	sessionRef := s.client.Doc(fmt.Sprintf("Session/%s", sessionId))
	snap, err := sessionRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("lookup sessionId %s: %w", sessionId, err)
	}
	var session store.Session
	if err = snap.DataTo(&session); err != nil {
		return nil, fmt.Errorf("lookup session %s: %w", sessionId, err)
	}
	session.LastUpdated = snap.UpdateTime.Format("2006-01-02T15:04:05Z")
	return &session, nil
}

func (s *Store) ListSessions(context context.Context, offset int, limit int) ([]*store.Session, error) {
	var sessions []*store.Session
	iter := s.client.Collection("Session").OrderBy("Id", firestore.Asc).Offset(offset).Limit(limit).Documents(context)
	for {
		snap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("next session: %w", err)
		}
		var session store.Session
		if err = snap.DataTo(&session); err != nil {
			return nil, fmt.Errorf("map session: %w", err)
		}
		session.LastUpdated = snap.UpdateTime.Format("2006-01-02T15:04:05Z")
		sessions = append(sessions, &session)
	}
	if sessions == nil {
		sessions = make([]*store.Session, 0)
	}
	return sessions, nil
}
