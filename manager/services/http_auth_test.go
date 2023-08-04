package services_test

import (
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"net/http"
	"os"
	"testing"
)

func TestEnvTokenHttpAuthService(t *testing.T) {
	err := os.Setenv("TEST_ENV_VAR", "test")
	require.NoError(t, err)
	svc := services.NewEnvTokenHttpAuthService("TEST_ENV_VAR")
	r, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	require.NoError(t, err)
	err = svc.Authenticate(r)
	require.NoError(t, err)
	require.Equal(t, "Bearer test", r.Header.Get("Authorization"))
	require.Equal(t, "test", r.Header.Get("X-API-Key"))
}
