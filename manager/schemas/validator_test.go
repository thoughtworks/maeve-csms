package schemas_test

import (
	"encoding/json"
	"errors"
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twlabs/maeve-csms/manager/ocpp/ocpp201"
	"github.com/twlabs/maeve-csms/manager/schemas"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	req := &ocpp201.BootNotificationRequestJson{
		ChargingStation: ocpp201.ChargingStationType{
			Model:      "SuperTest",
			VendorName: "TW",
		},
		Reason: ocpp201.BootReasonEnumTypePowerUp,
	}
	data, err := json.Marshal(req)
	require.NoError(t, err)

	err = schemas.Validate(data, schemas.OcppSchemas, "ocpp201/BootNotificationRequest.json")
	assert.NoError(t, err)
	var valError *jsonschema.ValidationError
	if errors.As(err, &valError) {
		t.Log(valError.Message)
	}
}

func TestValidateInvalidRequest(t *testing.T) {
	req := &ocpp201.BootNotificationRequestJson{
		ChargingStation: ocpp201.ChargingStationType{
			Model:      "SuperTest",
			VendorName: "TW",
		},
	}
	data, err := json.Marshal(req)
	require.NoError(t, err)

	err = schemas.Validate(data, schemas.OcppSchemas, "ocpp201/BootNotificationRequest.json")
	var valError *jsonschema.ValidationError
	assert.ErrorAs(t, err, &valError)
}
