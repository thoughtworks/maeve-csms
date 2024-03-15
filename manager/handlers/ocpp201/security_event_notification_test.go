// SPDX-License-Identifier: Apache-2.0

package ocpp201

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/testutil"
	"testing"
	"time"
)

func TestSecurityEventNotificationHandler(t *testing.T) {
	handler := SecurityEventNotificationHandler{}

	now := time.Now().UTC().Format(time.RFC3339)

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &ocpp201.SecurityEventNotificationRequestJson{
			Timestamp: now,
			Type:      "SomeSecurityEvent",
		}

		resp, err := handler.HandleCall(ctx, "cs001", req)
		require.NoError(t, err)

		assert.Equal(t, &ocpp201.SecurityEventNotificationResponseJson{}, resp)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"security_event.timestamp": now,
		"security_event.type":      "SomeSecurityEvent",
	})
}
