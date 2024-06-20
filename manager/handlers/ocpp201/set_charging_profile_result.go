// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"golang.org/x/exp/slog"
)

type SetChargingProfileResultHandler struct{}

func (h SetChargingProfileResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	req := request.(*types.SetChargingProfileRequestJson)
	resp := response.(*types.SetChargingProfileResponseJson)

	slog.Debug("[API TRACE] in scp_result.go, got response:", resp, "[API TRACE] From request:", req)

	return nil
}
