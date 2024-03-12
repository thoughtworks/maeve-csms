// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"testing"
)

func TestFirmwareStatusNotification(t *testing.T) {
	handler := ocpp201.FirmwareStatusNotificationHandler{}

	tracer, exporter := handlers.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		req := &types.FirmwareStatusNotificationRequestJson{
			Status: types.FirmwareStatusEnumTypeDownloading,
		}

		resp, err := handler.HandleCall(ctx, "cs001", req)
		require.NoError(t, err)

		assert.Equal(t, &types.FirmwareStatusNotificationResponseJson{}, resp)
	}()

	handlers.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"firmware_status.status": "Downloading",
	})
}

func TestFirmwareStatusNotificationWithRequestId(t *testing.T) {
	handler := ocpp201.FirmwareStatusNotificationHandler{}

	tracer, exporter := handlers.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		requestId := 42
		req := &types.FirmwareStatusNotificationRequestJson{
			Status:    types.FirmwareStatusEnumTypeDownloading,
			RequestId: &requestId,
		}

		resp, err := handler.HandleCall(ctx, "cs001", req)
		require.NoError(t, err)

		assert.Equal(t, &types.FirmwareStatusNotificationResponseJson{}, resp)
	}()

	handlers.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"firmware_status.status":     "Downloading",
		"firmware_status.request_id": 42,
	})
}
