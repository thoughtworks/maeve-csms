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

func TestChangeAvailabilityResultHandler(t *testing.T) {
	handler := ocpp201.ChangeAvailabilityResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.ChangeAvailabilityRequestJson{
			OperationalStatus: types.OperationalStatusEnumTypeOperative,
		}
		resp := &types.ChangeAvailabilityResponseJson{
			Status: types.ChangeAvailabilityStatusEnumTypeAccepted,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"change_availability.operational_status": "Operative",
		"change_availability.status":             "Accepted",
	})
}

func TestChangeAvailabilityResultHandlerWithEvseId(t *testing.T) {
	handler := ocpp201.ChangeAvailabilityResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.ChangeAvailabilityRequestJson{
			OperationalStatus: types.OperationalStatusEnumTypeOperative,
			Evse: &types.EVSEType{
				Id: 1,
			},
		}
		resp := &types.ChangeAvailabilityResponseJson{
			Status: types.ChangeAvailabilityStatusEnumTypeAccepted,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"change_availability.operational_status": "Operative",
		"change_availability.status":             "Accepted",
		"change_availability.evse_id":            1,
	})
}

func TestChangeAvailabilityResultHandlerWithConnectorId(t *testing.T) {
	handler := ocpp201.ChangeAvailabilityResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.ChangeAvailabilityRequestJson{
			OperationalStatus: types.OperationalStatusEnumTypeOperative,
			Evse: &types.EVSEType{
				Id:          1,
				ConnectorId: makePtr(2),
			},
		}
		resp := &types.ChangeAvailabilityResponseJson{
			Status: types.ChangeAvailabilityStatusEnumTypeAccepted,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"change_availability.operational_status": "Operative",
		"change_availability.status":             "Accepted",
		"change_availability.evse_id":            1,
		"change_availability.connector_id":       2,
	})
}
