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

func TestBaseReportResultHandler(t *testing.T) {
	handler := ocpp201.GetBaseReportResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.GetBaseReportRequestJson{
			RequestId:  42,
			ReportBase: types.ReportBaseEnumTypeFullInventory,
		}
		resp := &types.GetBaseReportResponseJson{
			Status: types.GenericDeviceModelStatusEnumTypeAccepted,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"get_base_report.request_id":  42,
		"get_base_report.report_base": "FullInventory",
		"get_base_report.status":      "Accepted",
	})
}
