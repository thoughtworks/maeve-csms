package firestore

import (
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
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
	return nil, nil
}

func (s *Store) ListSessions(context context.Context, offset int, limit int) ([]*store.Session, error) {
	return nil, nil
}
