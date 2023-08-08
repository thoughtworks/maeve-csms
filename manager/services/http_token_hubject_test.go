// SPDX-License-Identifier: Apache-2.0

//go:build integration

package services_test

import (
	"context"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"net/http"
	"testing"
)

func TestHubjectTestHttpTokenService(t *testing.T) {
	svc := services.NewHubjectTestHttpTokenService(
		"https://hubject.stoplight.io/api/v1/projects/cHJqOjk0NTg5/nodes/6bb8b3bc79c2e-authorization-token",
		http.DefaultClient)
	token, err := svc.GetToken(context.Background(), false)
	require.NoError(t, err)
	jwtToken, err := jwt.Parse([]byte(token))
	require.NoError(t, err)

	assert.Equal(t, "https://auth.eu.plugncharge.hubject.com/", jwtToken.Issuer())
}
