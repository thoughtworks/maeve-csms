// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"time"

	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"k8s.io/utils/clock"
)

type HeartbeatHandler struct {
	Clock clock.PassiveClock
}

func (h HeartbeatHandler) HandleCall(ctx context.Context, chargeStationId string, request ocpp.Request) (ocpp.Response, error) {
	return &types.HeartbeatResponseJson{
		CurrentTime: h.Clock.Now().Format(time.RFC3339),
	}, nil
}
