// SPDX-License-Identifier: Apache-2.0

package services_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"k8s.io/utils/clock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestEnvTokenService(t *testing.T) {
	err := os.Setenv("TEST_ENV_VAR", "test")
	require.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TEST_ENV_VAR")
	}()
	svc, err := services.NewEnvHttpTokenService("TEST_ENV_VAR")
	require.NoError(t, err)

	token, err := svc.GetToken(context.TODO(), false)
	require.NoError(t, err)
	assert.Equal(t, "test", token)
}

func TestEnvTokenServiceWithBadEnvVar(t *testing.T) {
	_, err := services.NewEnvHttpTokenService("BAD_ENV_VAR")
	require.Error(t, err)
}

func TestFixedTokenService(t *testing.T) {
	svc := services.NewFixedHttpTokenService("test")

	token, err := svc.GetToken(context.TODO(), false)
	require.NoError(t, err)
	assert.Equal(t, "test", token)
}

type CountingTokenService struct {
	Count int
	Token string
}

func (c *CountingTokenService) GetToken(_ context.Context, _ bool) (string, error) {
	c.Count++
	return fmt.Sprintf("%s%d", c.Token, c.Count), nil
}

func TestCachingTokenService(t *testing.T) {
	svc := services.NewCachingHttpTokenService(&CountingTokenService{Token: "test"}, 100*time.Millisecond, clock.RealClock{})
	token, err := svc.GetToken(context.Background(), false)
	require.NoError(t, err)
	assert.Equal(t, "test1", token)

	token, err = svc.GetToken(context.Background(), false)
	require.NoError(t, err)
	assert.Equal(t, "test1", token)

	time.Sleep(100 * time.Millisecond)

	token, err = svc.GetToken(context.Background(), false)
	require.NoError(t, err)
	assert.Equal(t, "test2", token)

	token, err = svc.GetToken(context.Background(), true)
	require.NoError(t, err)
	assert.Equal(t, "test3", token)
}

func TestOAuth2HttpTokenService(t *testing.T) {
	clientId := "client_id"
	clientSecret := "client_secret"

	handler := &oauth2HttpHandler{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()

	clk := clock.RealClock{}
	svc := services.NewOAuth2HttpTokenService(srv.URL, clientId, clientSecret, http.DefaultClient, clk)
	token, err := svc.GetToken(context.Background(), false)
	require.NoError(t, err)
	assert.Equal(t, "test", token)
}

func TestOAuth2HttpTokenServiceWithCachedCredential(t *testing.T) {
	clientId := "client_id"
	clientSecret := "client_secret"

	handler := &oauth2HttpHandler{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()

	clk := clock.RealClock{}
	svc := services.NewOAuth2HttpTokenService(srv.URL, clientId, clientSecret, http.DefaultClient, clk)
	token, err := svc.GetToken(context.Background(), false)
	require.NoError(t, err)
	assert.Equal(t, "test", token)

	token, err = svc.GetToken(context.Background(), false)
	require.NoError(t, err)
	assert.Equal(t, "test", token)

	assert.Equal(t, 1, handler.CallCount)
}

type oauth2HttpHandler struct {
	ClientId     string
	ClientSecret string
	CallCount    int
}

func (o *oauth2HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var tokenReq services.OAuth2TokenRequest
	err := json.NewDecoder(r.Body).Decode(&tokenReq)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if tokenReq.ClientId != o.ClientId || tokenReq.ClientSecret != o.ClientSecret {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else if tokenReq.GrantType != "client_credentials" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(services.OAuth2TokenResponse{
		AccessToken: "test",
		ExpiresIn:   10,
		TokenType:   "bearer",
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	o.CallCount++
}
