package ocpp201_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/testutil"
	"testing"
)

func TestGetVariablesResult(t *testing.T) {
	handler := ocpp201.GetVariablesResultHandler{}

	tracer, exporter := testutil.GetTracer()

	ctx := context.Background()

	func() {
		ctx, span := tracer.Start(ctx, `test`)
		defer span.End()

		req := &types.GetVariablesRequestJson{
			GetVariableData: []types.GetVariableDataType{
				{
					Component: types.ComponentType{
						Name: "SomeCtrlr",
					},
					Variable: types.VariableType{
						Name: "MyVar",
					},
					AttributeType: makePtr(types.AttributeEnumTypeMaxSet),
				},
				{
					Component: types.ComponentType{
						Name:     "SomeOtherCtrlr",
						Instance: makePtr("SomeInstance"),
					},
					Variable: types.VariableType{
						Name:     "MyOtherVar",
						Instance: makePtr("SomeVarInstance"),
					},
				},
				{
					Component: types.ComponentType{
						Name:     "SomeOtherCtrlr",
						Instance: makePtr("SomeInstance"),
						Evse: &types.EVSEType{
							Id: 1,
						},
					},
					Variable: types.VariableType{
						Name:     "MyOtherVar",
						Instance: makePtr("SomeVarInstance"),
					},
				},
				{
					Component: types.ComponentType{
						Name: "SomeCtrlr",
						Evse: &types.EVSEType{
							Id:          1,
							ConnectorId: makePtr(2),
						},
					},
					Variable: types.VariableType{
						Name: "AnotherVar",
					},
				},
			},
		}
		resp := &types.GetVariablesResponseJson{
			GetVariableResult: []types.GetVariableResultType{
				{
					Component: types.ComponentType{
						Name: "SomeCtrlr",
					},
					Variable: types.VariableType{
						Name: "MyVar",
					},
					AttributeType:   makePtr(types.AttributeEnumTypeMaxSet),
					AttributeValue:  makePtr("12"),
					AttributeStatus: types.GetVariableStatusEnumTypeAccepted,
				},
				{
					Component: types.ComponentType{
						Name:     "SomeOtherCtrlr",
						Instance: makePtr("SomeInstance"),
					},
					Variable: types.VariableType{
						Name:     "MyOtherVar",
						Instance: makePtr("SomeVarInstance"),
					},
					AttributeValue:  makePtr("Example"),
					AttributeStatus: types.GetVariableStatusEnumTypeAccepted,
				},
				{
					Component: types.ComponentType{
						Name:     "SomeOtherCtrlr",
						Instance: makePtr("SomeInstance"),
						Evse: &types.EVSEType{
							Id: 1,
						},
					},
					Variable: types.VariableType{
						Name:     "MyOtherVar",
						Instance: makePtr("SomeVarInstance"),
					},
					AttributeValue:  nil,
					AttributeStatus: types.GetVariableStatusEnumTypeNotSupportedAttributeType,
				},
				{
					Component: types.ComponentType{
						Name: "SomeCtrlr",
						Evse: &types.EVSEType{
							Id:          1,
							ConnectorId: makePtr(2),
						},
					},
					Variable: types.VariableType{
						Name: "AnotherVar",
					},
					AttributeValue:  makePtr("Hello"),
					AttributeStatus: types.GetVariableStatusEnumTypeAccepted,
				},
			},
		}

		err := handler.HandleCallResult(ctx, "cs001", req, resp, nil)
		require.NoError(t, err)
	}()

	testutil.AssertSpan(t, &exporter.GetSpans()[0], "test", map[string]any{
		"get_variables.names":  "SomeCtrlr:::/MyVar:,SomeOtherCtrlr:SomeInstance::/MyOtherVar:SomeVarInstance,SomeOtherCtrlr:SomeInstance:1:/MyOtherVar:SomeVarInstance,SomeCtrlr::1:2/AnotherVar:",
		"get_variables.values": "MaxSet/12:Accepted,Actual/Example:Accepted,Actual/<null>:NotSupportedAttributeType,Actual/Hello:Accepted",
	})
}
