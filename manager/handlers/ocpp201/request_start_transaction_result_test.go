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

func TestRequestStartTransactionResultHandler(t *testing.T) {
	handler := ocpp201.RequestStartTransactionResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.RequestStartTransactionRequestJson{
			RemoteStartId: 123,
			IdToken: types.IdTokenType{
				Type:    types.IdTokenEnumTypeISO14443,
				IdToken: "DEADBEEF",
			},
		}
		resp := &types.RequestStartTransactionResponseJson{
			Status: types.RequestStartStopStatusEnumTypeAccepted,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"request_start.remote_start_id": 123,
		"request_start.status":          "Accepted",
	})
}

func TestRequestStartTransactionResultHandlerWithTransactionId(t *testing.T) {
	handler := ocpp201.RequestStartTransactionResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.RequestStartTransactionRequestJson{
			RemoteStartId: 123,
			IdToken: types.IdTokenType{
				Type:    types.IdTokenEnumTypeISO14443,
				IdToken: "DEADBEEF",
			},
		}
		resp := &types.RequestStartTransactionResponseJson{
			Status:        types.RequestStartStopStatusEnumTypeAccepted,
			TransactionId: makePtr("abc12345"),
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"request_start.remote_start_id": 123,
		"request_start.status":          "Accepted",
		"request_start.transaction_id":  "abc12345",
	})
}
