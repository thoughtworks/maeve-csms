// SPDX-License-Identifier: Apache-2.0

package ocpi_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpi"
	"github.com/thoughtworks/maeve-csms/manager/server"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistration(t *testing.T) {
	tokenA := "abcdef123456"

	// setup sender
	senderStore := inmemory.NewStore(clock.RealClock{})
	senderOcpiApi := ocpi.NewOCPI(senderStore, http.DefaultClient, "GB", "TWK")
	senderHandler := server.NewOcpiHandler(senderStore, clock.RealClock{}, senderOcpiApi, newNoopV16CallMaker())
	senderServer := httptest.NewServer(senderHandler)
	senderOcpiApi.SetExternalUrl(senderServer.URL)
	defer senderServer.Close()

	// setup receiver
	receiverStore := inmemory.NewStore(clock.RealClock{})
	err := receiverStore.SetRegistrationDetails(context.Background(), tokenA, &store.OcpiRegistration{
		Status: store.OcpiRegistrationStatusPending,
	})
	require.NoError(t, err)
	receiverOcpiApi := ocpi.NewOCPI(receiverStore, http.DefaultClient, "GB", "TWS")
	receiverHandler := server.NewOcpiHandler(receiverStore, clock.RealClock{}, receiverOcpiApi, newNoopV16CallMaker())
	receiverServer := httptest.NewServer(receiverHandler)
	receiverOcpiApi.SetExternalUrl(receiverServer.URL)
	defer receiverServer.Close()

	// registration
	err = senderOcpiApi.RegisterNewParty(context.Background(), receiverServer.URL+"/ocpi/versions", tokenA)
	require.NoError(t, err)

	// check registration status
	senderPartyDetails, err := receiverStore.GetPartyDetails(context.Background(), "CPO", "GB", "TWK")
	require.NoError(t, err)

	assert.Equal(t, "CPO", senderPartyDetails.Role)
	assert.Equal(t, "GB", senderPartyDetails.CountryCode)
	assert.Equal(t, "TWK", senderPartyDetails.PartyId)
	assert.Equal(t, senderServer.URL+"/ocpi/versions", senderPartyDetails.Url)
	assert.Len(t, senderPartyDetails.Token, 64)

	receiverPartyDetails, err := senderStore.GetPartyDetails(context.Background(), "CPO", "GB", "TWS")
	require.NoError(t, err)

	assert.Equal(t, "CPO", receiverPartyDetails.Role)
	assert.Equal(t, "GB", receiverPartyDetails.CountryCode)
	assert.Equal(t, "TWS", receiverPartyDetails.PartyId)
	assert.Equal(t, receiverServer.URL+"/ocpi/versions", receiverPartyDetails.Url)
	assert.Len(t, receiverPartyDetails.Token, 64)

	// check initial registration details have been removed
	senderTokenAReg, err := senderStore.GetRegistrationDetails(context.Background(), tokenA)
	require.NoError(t, err)
	assert.Nil(t, senderTokenAReg)

	receiverTokenAReg, err := receiverStore.GetRegistrationDetails(context.Background(), tokenA)
	require.NoError(t, err)
	assert.Nil(t, receiverTokenAReg)

	// token status
	receiverTokenBReg, err := receiverStore.GetRegistrationDetails(context.Background(), senderPartyDetails.Token)
	require.NoError(t, err)
	require.NotNil(t, receiverTokenBReg)
	assert.Equal(t, store.OcpiRegistrationStatusRegistered, receiverTokenBReg.Status)

	senderTokenCReg, err := senderStore.GetRegistrationDetails(context.Background(), receiverPartyDetails.Token)
	require.NoError(t, err)
	require.NotNil(t, senderTokenCReg)
	assert.Equal(t, store.OcpiRegistrationStatusRegistered, senderTokenCReg.Status)
}
