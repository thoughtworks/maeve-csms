// SPDX-License-Identifier: Apache-2.0

package ocpp16

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	"github.com/thoughtworks/maeve-csms/manager/store"
)

type MeterValuesHandler struct {
	TransactionStore store.TransactionStore
}

func (m MeterValuesHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (response ocpp.Response, err error) {
	// TODO: store in transaction store

	return &types.MeterValuesResponseJson{}, nil
}
