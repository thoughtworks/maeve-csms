package bigtable

import (
	"bytes"
	bt "cloud.google.com/go/bigtable"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
)

const (
	csType                = "cs"
	securityProfileColumn = "prof"
	passwordColumn        = "pwd"
)

func (c *Store) SetChargeStationAuth(ctx context.Context, chargeStationId string, auth *store.ChargeStationAuth) error {
	mut := bt.NewMutation()

	secProfile := new(bytes.Buffer)
	err := binary.Write(secProfile, binary.BigEndian, auth.SecurityProfile)
	if err != nil {
		return err
	}

	mut.Set(idColumnFamily, typeColumn, bt.Now(), []byte(csType))
	mut.Set(idColumnFamily, csType, bt.Now(), []byte(chargeStationId))
	mut.Set(csAuthColumnFamily, securityProfileColumn, bt.Now(), secProfile.Bytes())
	mut.Set(csAuthColumnFamily, passwordColumn, bt.Now(), []byte(auth.Base64SHA256Password))

	rowKey := fmt.Sprintf("%s#%s", csType, chargeStationId)

	return c.table.Apply(ctx, rowKey, mut)
}

func (c *Store) LookupChargeStationAuth(ctx context.Context, chargeStationId string) (*store.ChargeStationAuth, error) {
	rowKey := fmt.Sprintf("%s#%s", csType, chargeStationId)

	row, err := c.table.ReadRow(ctx, rowKey, bt.RowFilter(bt.FamilyFilter(csAuthColumnFamily)))
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}
	auth := new(store.ChargeStationAuth)
	for columnFamily, column := range row {
		if columnFamily == csAuthColumnFamily {
			for _, col := range column {
				switch col.Column {
				case fmt.Sprintf("%s:%s", csAuthColumnFamily, securityProfileColumn):
					err = binary.Read(bytes.NewReader(col.Value), binary.BigEndian, &auth.SecurityProfile)
					if err != nil {
						return nil, fmt.Errorf("read security profile: %w", err)
					}
				case fmt.Sprintf("%s:%s", csAuthColumnFamily, passwordColumn):
					auth.Base64SHA256Password = string(col.Value)
				}
			}
		}
	}

	return auth, nil
}
