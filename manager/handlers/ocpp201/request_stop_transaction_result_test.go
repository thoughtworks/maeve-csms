// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/testutil"
	"testing"
)

func TestRequestStopTransactionResultHandler(t *testing.T) {
	handler := ocpp201.RequestStopTransactionResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.RequestStopTransactionRequestJson{
			TransactionId: "abc12345",
		}
		resp := &types.RequestStopTransactionResponseJson{
			Status: types.RequestStartStopStatusEnumTypeAccepted,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"request_stop.transaction_id": "abc12345",
		"request_stop.status":         "Accepted",
	})
}
