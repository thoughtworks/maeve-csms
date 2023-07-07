// SPDX-License-Identifier: Apache-2.0

package services_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"testing"
)

func TestInMemoryTokenStore(t *testing.T) {
	emaidToken := &services.Token{
		Type: string(ocpp201.IdTokenEnumTypeEMAID),
		Uid:  "SOMEEMAID",
	}
	rfidToken := &services.Token{
		Type: string(ocpp201.IdTokenEnumTypeISO14443),
		Uid:  "SOMERFID",
	}

	tokenStore := services.InMemoryTokenStore{
		Tokens: map[string]*services.Token{
			"eMAID:SOMEEMAID":   emaidToken,
			"ISO14443:SOMERFID": rfidToken,
		},
	}

	token, err := tokenStore.FindToken(string(ocpp201.IdTokenEnumTypeEMAID), "SOMEEMAID")
	require.NoError(t, err)
	assert.Equal(t, emaidToken, token)

	token, err = tokenStore.FindToken(string(ocpp201.IdTokenEnumTypeISO14443), "SOMERFID")
	require.NoError(t, err)
	assert.Equal(t, rfidToken, token)

	token, err = tokenStore.FindToken(string(ocpp201.IdTokenEnumTypeISO14443), "SOMEEMAID")
	require.NoError(t, err)
	assert.Nil(t, token)

	token, err = tokenStore.FindToken("", "SOMERFID")
	require.NoError(t, err)
	assert.Equal(t, rfidToken, token)

	token, err = tokenStore.FindToken("", "SOMEEMAID")
	require.NoError(t, err)
	assert.Equal(t, emaidToken, token)
}
