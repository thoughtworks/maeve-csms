package ocpp201_test

import (
	"context"
	"github.com/stretchr/testify/require"
	handlers201 "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	"testing"
)

func TestSetVariablesResultHandler(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	handler := handlers201.SetVariablesResultHandler{
		Store: engine,
	}

	request := ocpp201.SetVariablesRequestJson{
		SetVariableData: []ocpp201.SetVariableDataType{
			{
				AttributeValue: "20",
				Component: ocpp201.ComponentType{
					Name: "MyCtrlr",
				},
				Variable: ocpp201.VariableType{
					Name: "MyVariable",
				},
			},
			{
				AttributeValue: "something",
				Component: ocpp201.ComponentType{
					Name: "MyCtrlr",
				},
				Variable: ocpp201.VariableType{
					Name: "MyOtherVariable",
				},
			},
		},
	}

	response := ocpp201.SetVariablesResponseJson{
		SetVariableResult: []ocpp201.SetVariableResultType{
			{
				Component: ocpp201.ComponentType{
					Name: "MyCtrlr",
				},
				Variable: ocpp201.VariableType{
					Name: "MyVariable",
				},
				AttributeStatus: ocpp201.SetVariableStatusEnumTypeAccepted,
			},
			{
				Component: ocpp201.ComponentType{
					Name: "MyCtrlr",
				},
				Variable: ocpp201.VariableType{
					Name: "MyOtherVariable",
				},
				AttributeStatus: ocpp201.SetVariableStatusEnumTypeRejected,
			},
		},
	}

	err := handler.HandleCallResult(context.TODO(), "cs001", &request, &response, nil)
	require.NoError(t, err)
}
