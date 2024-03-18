// SPDX-License-Identifier: Apache-2.0

package ocpp201_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/testutil"
	"testing"
)

func TestNotifyReport(t *testing.T) {
	handler := ocpp201.NotifyReportHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		req := &types.NotifyReportRequestJson{
			GeneratedAt: "2024-03-18T17:10:00.000Z",
			ReportData: []types.ReportDataType{
				{
					Component: types.ComponentType{
						Name: "SomeCtrlr",
					},
					Variable: types.VariableType{
						Name: "SomeVar",
					},
					VariableAttribute: []types.VariableAttributeType{
						{
							Constant:   false,
							Mutability: makePtr(types.MutabilityEnumTypeReadOnly),
							Persistent: true,
							Type:       makePtr(types.AttributeEnumTypeActual),
							Value:      makePtr("19"),
						},
					},
					VariableCharacteristics: nil,
				},
			},
			RequestId: 42,
			SeqNo:     1,
			Tbc:       false,
		}

		resp, err := handler.HandleCall(ctx, "cs001", req)
		require.NoError(t, err)

		assert.Equal(t, &types.NotifyReportResponseJson{}, resp)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"notify_report.generated_at": "2024-03-18T17:10:00.000Z",
		"notify_report.request_id":   42,
		"notify_report.seq_no":       1,
		"notify_report.tbc":          false,
	})
}
