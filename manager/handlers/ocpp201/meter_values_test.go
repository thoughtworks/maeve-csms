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

func TestMeterValuesHandler(t *testing.T) {
	handler := ocpp201.MeterValuesHandler{}

	tracer, exporter := handlers.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, "test")
		defer span.End()

		req := &types.MeterValuesRequestJson{
			EvseId: 1,
			MeterValue: []types.MeterValueType{
				{
					SampledValue: []types.SampledValueType{
						{
							Measurand: makePtr(types.MeasurandEnumTypeEnergyActiveImportRegister),
							Location:  makePtr(types.LocationEnumTypeOutlet),
							Value:     100,
						},
					},
					Timestamp: "2023-06-15T15:05:00+01:00",
				},
			},
		}

		resp, err := handler.HandleCall(ctx, "cs001", req)
		require.NoError(t, err)

		assert.Equal(t, &types.MeterValuesResponseJson{}, resp)
	}()

	handlers.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"meter_values.evse_id": 1,
	})

}
