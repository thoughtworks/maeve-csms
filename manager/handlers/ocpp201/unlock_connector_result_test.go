// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"testing"
)

func TestUnlockConnectorResult(t *testing.T) {
	handler := ocpp201.UnlockConnectorResultHandler{}

	tracer, exporter := handlers.GetTracer()

	ctx := context.TODO()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.UnlockConnectorRequestJson{
			EvseId:      1,
			ConnectorId: 2,
		}

		resp := &types.UnlockConnectorResponseJson{
			Status: types.UnlockStatusEnumTypeUnlocked,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	handlers.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"unlock_connector.evse_id":      1,
		"unlock_connector.connector_id": 2,
		"unlock_connector.status":       "Unlocked",
	})
}
