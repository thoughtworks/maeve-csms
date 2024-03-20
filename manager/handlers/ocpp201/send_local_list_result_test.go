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

func TestSendLocalListResultHandler(t *testing.T) {
	handler := ocpp201.SendLocalListResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.SendLocalListRequestJson{
			LocalAuthorizationList: []types.AuthorizationData{
				{
					IdToken: types.IdTokenType{
						Type:    types.IdTokenEnumTypeISO14443,
						IdToken: "ABCD1234",
					},
				},
			},
			UpdateType:    types.UpdateEnumTypeFull,
			VersionNumber: 42,
		}
		resp := &types.SendLocalListResponseJson{
			Status: types.SendLocalListStatusEnumTypeAccepted,
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"send_local_list.update_type":    "Full",
		"send_local_list.version_number": 42,
		"send_local_list.status":         "Accepted",
	})
}
