// SPDX-License-Identifier: Apache-2.0

package ocpp16_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp16"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp16"
	clockTest "k8s.io/utils/clock/testing"
	"testing"
	"time"
)

func TestBootNotificationHandler(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2023-06-15T15:05:00+01:00")
	require.NoError(t, err)
	clock := clockTest.NewFakePassiveClock(now)

	handler := handlers.BootNotificationHandler{
		Clock:             clock,
		HeartbeatInterval: 10,
	}

	serialNumber := "cs001-1234"
	req := &types.BootNotificationJson{
		ChargePointSerialNumber: &serialNumber,
	}

	got, err := handler.HandleCall(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.BootNotificationResponseJson{
		CurrentTime: "2023-06-15T15:05:00+01:00",
		Status:      types.BootNotificationResponseJsonStatusAccepted,
		Interval:    10,
	}

	assert.Equal(t, want, got)
}
