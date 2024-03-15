package ocpp201

import (
	"context"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	types "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

type GetVariablesResultHandler struct{}

func (h GetVariablesResultHandler) HandleCallResult(ctx context.Context, chargeStationId string, request ocpp.Request, response ocpp.Response, state any) error {
	resp := response.(*types.GetVariablesResponseJson)

	span := trace.SpanFromContext(ctx)

	var variableNames []string
	var variableValues []string
	for _, v := range resp.GetVariableResult {
		variableNames = append(variableNames, fmt.Sprintf("%s/%s", getComponentId(v.Component), getVariableId(v.Variable)))
		variableValues = append(variableValues, fmt.Sprintf("%s/%s:%s", getAttributeTypeName(v.AttributeType), getAttributeValue(v.AttributeValue), v.AttributeStatus))
	}

	span.SetAttributes(
		attribute.String("get_variables.names", strings.Join(variableNames, ",")),
		attribute.String("get_variables.values", strings.Join(variableValues, ",")))

	return nil
}

func getComponentId(component types.ComponentType) string {
	instance := ""
	evseId := ""
	connectorId := ""

	if component.Instance != nil {
		instance = *component.Instance
	}
	if component.Evse != nil {
		evseId = fmt.Sprintf("%d", component.Evse.Id)
		if component.Evse.ConnectorId != nil {
			connectorId = fmt.Sprintf("%d", *component.Evse.ConnectorId)
		}
	}

	return fmt.Sprintf("%s:%s:%s:%s", component.Name, instance, evseId, connectorId)
}

func getVariableId(variable types.VariableType) string {
	instance := ""

	if variable.Instance != nil {
		instance = *variable.Instance
	}

	return fmt.Sprintf("%s:%s", variable.Name, instance)
}

func getAttributeTypeName(attributeType *types.AttributeEnumType) string {
	result := "Actual"
	if attributeType != nil {
		result = string(*attributeType)
	}
	return result
}

func getAttributeValue(value *string) string {
	result := "<null>"
	if value != nil {
		result = *value
	}
	return result
}
