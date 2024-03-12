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

func TestLogStatusNotification(t *testing.T) {
	handler := ocpp201.LogStatusNotificationHandler{}

	tracer, exporter := handlers.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		req := &types.LogStatusNotificationRequestJson{
			Status: types.UploadLogStatusEnumTypeUploaded,
		}

		resp, err := handler.HandleCall(ctx, "cs001", req)
		require.NoError(t, err)

		assert.Equal(t, &types.LogStatusNotificationResponseJson{}, resp)
	}()

	handlers.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"log_status.status": "Uploaded",
	})
}

func TestLogStatusNotificationWithRequestId(t *testing.T) {
	handler := ocpp201.LogStatusNotificationHandler{}

	tracer, exporter := handlers.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		requestId := 999
		req := &types.LogStatusNotificationRequestJson{
			Status:    types.UploadLogStatusEnumTypeIdle,
			RequestId: &requestId,
		}

		resp, err := handler.HandleCall(ctx, "cs001", req)
		require.NoError(t, err)

		assert.Equal(t, &types.LogStatusNotificationResponseJson{}, resp)
	}()

	handlers.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"log_status.status":     "Idle",
		"log_status.request_id": 999,
	})
}
