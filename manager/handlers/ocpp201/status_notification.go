// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"golang.org/x/exp/slog"
)

func StatusNotificationHandler(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	req := request.(*types.StatusNotificationRequestJson)
	slog.Info("status notification", slog.String("chargeStationId", chargeStationId),
		slog.Int("evseId", req.EvseId), slog.Int("connectorId", req.ConnectorId), slog.Any("connectorStatus", req.ConnectorStatus))
	return &types.StatusNotificationResponseJson{}, nil
}
