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

func TestTransactionStatusResultHandlerWithTransactionId(t *testing.T) {
	handler := ocpp201.GetTransactionStatusResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.GetTransactionStatusRequestJson{
			TransactionId: makePtr("1234567890"),
		}
		resp := &types.GetTransactionStatusResponseJson{
			MessagesInQueue: true,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"get_transaction_status.transaction_id":    "1234567890",
		"get_transaction_status.messages_in_queue": true,
	})
}

func TestTransactionStatusResultHandlerWithOngoingIndicator(t *testing.T) {
	handler := ocpp201.GetTransactionStatusResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.GetTransactionStatusRequestJson{}
		resp := &types.GetTransactionStatusResponseJson{
			MessagesInQueue:  false,
			OngoingIndicator: makePtr(true),
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"get_transaction_status.messages_in_queue": false,
		"get_transaction_status.ongoing":           true,
	})
}
