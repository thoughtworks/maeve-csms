package firestore

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/store"
)

func (s *Store) SetSession(ctx context.Context, session *store.Session) error {

	return nil
}

func (s *Store) LookupSession(ctx context.Context, sessionId string) (*store.Session, error) {
	return nil, nil
}

func (s *Store) ListSessions(context context.Context, offset int, limit int) ([]*store.Session, error) {
	return nil, nil
}
