// SPDX-License-Identifier: Apache-2.0

package ocpi_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/ocpi"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"github.com/thoughtworks/maeve-csms/manager/store/inmemory"
	"k8s.io/utils/clock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetVersions(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")

	want := []ocpi.Version{
		{
			Version: "2.2",
			Url:     "/ocpi/2.2",
		},
	}

	got, err := ocpiApi.GetVersions(context.Background())
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetVersionDetails(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")

	want := ocpi.VersionDetail{
		Version: "2.2",
		Endpoints: []ocpi.Endpoint{
			{
				Identifier: "credentials",
				Role:       ocpi.RECEIVER,
				Url:        "/ocpi/2.2/credentials",
			},
			{
				Identifier: "commands",
				Role:       ocpi.RECEIVER,
				Url:        "/ocpi/receiver/2.2/commands",
			},
			{
				Identifier: "tokens",
				Role:       ocpi.RECEIVER,
				Url:        "/ocpi/receiver/2.2/tokens/",
			},
		},
	}

	got, err := ocpiApi.GetVersion(context.Background())
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetToken(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")
	err := engine.SetToken(context.Background(), &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		ContractId:  "GBTWKTWTW000018",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "ALWAYS",
	})
	require.NoError(t, err)

	want := &ocpi.Token{
		ContractId:  "GBTWKTWTW000018",
		CountryCode: "GB",
		Issuer:      "Thoughtworks",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		Valid:       true,
		Whitelist:   "ALWAYS",
	}

	got, err := ocpiApi.GetToken(context.Background(), "GB", "TWK", "DEADBEEF")
	require.NoError(t, err)

	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`, got.LastUpdated)
	got.LastUpdated = ""
	assert.Equal(t, want, got)
}

func TestSetToken(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")

	err := ocpiApi.SetToken(context.Background(), ocpi.Token{
		ContractId:  "GBTWKTWTW000018",
		CountryCode: "GB",
		Issuer:      "Thoughtworks",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		Valid:       true,
		Whitelist:   "ALWAYS",
	})
	require.NoError(t, err)

	want := &store.Token{
		CountryCode: "GB",
		PartyId:     "TWK",
		Type:        "RFID",
		Uid:         "DEADBEEF",
		ContractId:  "GBTWKTWTW000018",
		Issuer:      "Thoughtworks",
		Valid:       true,
		CacheMode:   "ALWAYS",
	}

	got, err := engine.LookupToken(context.Background(), "DEADBEEF")
	require.NoError(t, err)
	got.LastUpdated = ""
	assert.Equal(t, want, got)
}

func TestPushLocation(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")

	mux := http.NewServeMux()
	receiverServer := httptest.NewServer(mux)
	defer receiverServer.Close()
	mux.HandleFunc("/ocpi/versions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"data":[{"version":"2.2","url":"%s/ocpi/2.2"}], "status_code":1000}`, receiverServer.URL)))
	})
	mux.HandleFunc("/ocpi/2.2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"data":{
				"version":"2.2",
				"endpoints":[{"identifier":"locations","role":"RECEIVER","url":"%s/ocpi/receiver/2.2/locations"}]},
				"status_code":1000}`,
			receiverServer.URL)))
	})
	mux.HandleFunc("/ocpi/receiver/2.2/locations/GB/TWK/loc001", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, r.Header.Get("X-Correlation-ID"), "some-correlation-id")
		var got []byte
		r.Body.Read(got)
		assert.Equal(t, []byte(nil), got)
		w.WriteHeader(http.StatusCreated)
	})
	err := ocpiApi.SetCredentials(context.Background(), "some-token-123", ocpi.Credentials{
		Roles: []ocpi.CredentialsRole{
			{
				CountryCode: "GB",
				PartyId:     "TWK",
				Role:        ocpi.CredentialsRoleRoleEMSP,
			},
		},
		Token: "some-token-456",
		Url:   receiverServer.URL + "/ocpi/versions",
	})
	require.NoError(t, err)

	ctxWithCorrelationId := context.WithValue(context.Background(), ocpi.ContextKeyCorrelationId, "some-correlation-id")
	err = ocpiApi.PushLocation(ctxWithCorrelationId, ocpi.Location{Id: "loc001"})

	require.NoError(t, err)
}

func TestPushSession(t *testing.T) {
	engine := inmemory.NewStore(clock.RealClock{})
	ocpiApi := ocpi.NewOCPI(engine, http.DefaultClient, "GB", "TWK")
	mux := http.NewServeMux()
	receiverServer := httptest.NewServer(mux)
	defer receiverServer.Close()
	mux.HandleFunc("/ocpi/versions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"data":[{"version":"2.2","url":"%s/ocpi/2.2"}], "status_code":1000}`, receiverServer.URL)))
	})
	mux.HandleFunc("/ocpi/2.2", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"data":{
				"version":"2.2",
				"endpoints":[{"identifier":"locations","role":"RECEIVER","url":"%s/ocpi/receiver/2.2/locations"}]},
				"status_code":1000}`,
			receiverServer.URL)))
	})
	mux.HandleFunc("/ocpi/receiver/2.2/sessions/GB/TWK/s001", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, r.Header.Get("X-Correlation-ID"), "some-correlation-id")
		var got []byte
		r.Body.Read(got)
		assert.Equal(t, []byte(nil), got)
		w.WriteHeader(http.StatusCreated)
	})
	err := ocpiApi.SetCredentials(context.Background(), "some-token-123", ocpi.Credentials{
		Roles: []ocpi.CredentialsRole{
			{
				CountryCode: "GB",
				PartyId:     "TWK",
				Role:        ocpi.CredentialsRoleRoleEMSP,
			},
		},
		Token: "some-token-456",
		Url:   receiverServer.URL + "/ocpi/versions",
	})
	require.NoError(t, err)
	token, _ := engine.LookupToken(context.Background(), "some-token-123")
	err = ocpiApi.PushSession(context.Background(), *token)

	require.NoError(t, err)
}
