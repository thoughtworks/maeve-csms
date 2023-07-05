package ocpp16_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	handlers "github.com/twlabs/maeve-csms/manager/handlers/ocpp16"
	types "github.com/twlabs/maeve-csms/manager/ocpp/ocpp16"
	"testing"
)

func TestStatusNotificationHandler(t *testing.T) {
	timestamp := "2023-05-01T01:00:00+01:00"
	req := &types.StatusNotificationJson{
		Timestamp:   &timestamp,
		ConnectorId: 2,
		ErrorCode:   types.StatusNotificationJsonErrorCodeNoError,
		Status:      types.StatusNotificationJsonStatusPreparing,
	}

	got, err := handlers.StatusNotificationHandler(context.Background(), "cs001", req)
	assert.NoError(t, err)

	want := &types.StatusNotificationResponseJson{}

	assert.Equal(t, want, got)
}
