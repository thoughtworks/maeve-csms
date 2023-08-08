// SPDX-License-Identifier: Apache-2.0

package services_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/services"
	"k8s.io/utils/clock"
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
