package ocpp201_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	handlers "github.com/twlabs/ocpp2-broker-core/manager/handlers/ocpp201"
	types "github.com/twlabs/ocpp2-broker-core/manager/ocpp/ocpp201"
	"testing"
)

func TestStatusNotificationHandler(t *testing.T) {
	req := &types.StatusNotificationRequestJson{
		Timestamp:       "2023-05-01T01:00:00+01:00",
		EvseId:          1,
		ConnectorId:     2,
		ConnectorStatus: types.ConnectorStatusEnumTypeOccupied,
	}

	got, err := handlers.StatusNotificationHandler(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.StatusNotificationResponseJson{}

	assert.Equal(t, want, got)
}
