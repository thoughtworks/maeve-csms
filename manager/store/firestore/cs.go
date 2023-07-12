package firestore

import (
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type chargeStation struct {
	SecurityProfile      int    `firestore:"prof"`
	Base64SHA256Password string `firestore:"pwd"`
}

func (s *Store) SetChargeStationAuth(ctx context.Context, chargeStationId string, auth *store.ChargeStationAuth) error {
	csRef := s.client.Doc(fmt.Sprintf("ChargeStation/%s", chargeStationId))
	_, err := csRef.Set(ctx, &chargeStation{
		SecurityProfile:      int(auth.SecurityProfile),
		Base64SHA256Password: auth.Base64SHA256Password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) LookupChargeStationAuth(ctx context.Context, chargeStationId string) (*store.ChargeStationAuth, error) {
	csRef := s.client.Doc(fmt.Sprintf("ChargeStation/%s", chargeStationId))
	snap, err := csRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("lookup charge station %s: %w", chargeStationId, err)
	}
	var csData chargeStation
	if err = snap.DataTo(&csData); err != nil {
		return nil, fmt.Errorf("map charge station %s: %w", chargeStationId, err)
	}
	return &store.ChargeStationAuth{
		SecurityProfile:      store.SecurityProfile(csData.SecurityProfile),
		Base64SHA256Password: csData.Base64SHA256Password,
	}, nil
}
