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

func TestSetNetworkProfileResultHandler(t *testing.T) {
	handler := ocpp201.SetNetworkProfileResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.SetNetworkProfileRequestJson{
			ConfigurationSlot: 1,
			ConnectionData: types.NetworkConnectionProfileType{
				MessageTimeout:  30,
				OcppCsmsUrl:     "https://cs.example.com/",
				OcppInterface:   types.OCPPInterfaceEnumTypeWired0,
				OcppTransport:   types.OCPPTransportEnumTypeJSON,
				OcppVersion:     types.OCPPVersionEnumTypeOCPP20,
				SecurityProfile: 2,
			},
		}
		resp := &types.SetNetworkProfileResponseJson{
			Status: types.SetNetworkProfileStatusEnumTypeAccepted,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"set_network_profile.config_slot": 1,
		"set_network_profile.status":      "Accepted",
	})
}
