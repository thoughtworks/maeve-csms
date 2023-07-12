package inmemory

import (
	"context"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"sync"
)

// Store is an in-memory implementation of the store.Engine interface. As everything
// is stored in memory it is not stateless and cannot be used if running >1 manager
// instances. It is primarily provided to support unit testing.
type Store struct {
	sync.Mutex
	chargeStationAuth map[string]*store.ChargeStationAuth
	tokens            map[string]*store.Token
}

func NewStore() *Store {
	return &Store{
		chargeStationAuth: make(map[string]*store.ChargeStationAuth),
		tokens:            make(map[string]*store.Token),
	}
}

func (s *Store) SetChargeStationAuth(ctx context.Context, chargeStationId string, auth *store.ChargeStationAuth) error {
	s.Lock()
	defer s.Unlock()
	s.chargeStationAuth[chargeStationId] = auth
	return nil
}

func (s *Store) LookupChargeStationAuth(ctx context.Context, chargeStationId string) (*store.ChargeStationAuth, error) {
	s.Lock()
	defer s.Unlock()
	return s.chargeStationAuth[chargeStationId], nil
}

func (s *Store) SetToken(ctx context.Context, token *store.Token) error {
	s.Lock()
	defer s.Unlock()
	s.tokens[token.Uid] = token
	return nil
}

func (s *Store) LookupToken(ctx context.Context, tokenUid string) (*store.Token, error) {
	s.Lock()
	defer s.Unlock()
	return s.tokens[tokenUid], nil
}
