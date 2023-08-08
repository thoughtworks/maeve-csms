// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/config"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestConfigure(t *testing.T) {
	cfg := &config.DefaultConfig
	cfg.ContractCertValidator.Ocsp.RootCertProvider.File.FileNames = []string{"testdata/root_ca.pem"}

	settings, err := config.Configure(context.TODO(), cfg)
	require.NoError(t, err)

	brokerUrl, err := url.Parse("mqtt://localhost:1883")
	require.NoError(t, err)

	wantApiSettings := config.ApiSettings{
		Addr: "localhost:9410",
	}

	wantMqttSettings := config.MqttSettings{
		Urls:              []*url.URL{brokerUrl},
		Prefix:            "cs",
		Group:             "manager",
		ConnectTimeout:    10 * time.Second,
		ConnectRetryDelay: 1 * time.Second,
		KeepAliveInterval: 10 * time.Second,
	}

	assert.Equal(t, wantApiSettings, settings.Api)
	assert.Equal(t, wantMqttSettings, settings.Mqtt)
	assert.NotNil(t, settings.Storage)
	assert.NotNil(t, settings.ContractCertValidationService)
	assert.NotNil(t, settings.ContractCertProviderService)
	assert.NotNil(t, settings.ChargeStationCertProviderService)
	assert.NotNil(t, settings.TariffService)
	assert.NotNil(t, settings.TracerProvider)
}

func TestConfigureFirestoreStorage(t *testing.T) {
	_ = os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")

	cfg := &config.DefaultConfig
	cfg.ContractCertValidator.Ocsp.RootCertProvider.File.FileNames = []string{"testdata/root_ca.pem"}
	cfg.Storage.Type = "firestore"
	cfg.Storage.FirestoreStorage = &config.FirestoreStorageConfig{
		ProjectId: "test-project-id",
	}

	settings, err := config.Configure(context.TODO(), cfg)
	require.NoError(t, err)
	require.NotNil(t, settings.Storage)
}

func TestConfigureInMemoryStorage(t *testing.T) {
	cfg := &config.DefaultConfig
	cfg.ContractCertValidator.Ocsp.RootCertProvider.File.FileNames = []string{"testdata/root_ca.pem"}
	cfg.Storage.Type = "in_memory"

	settings, err := config.Configure(context.TODO(), cfg)
	require.NoError(t, err)
	require.NotNil(t, settings.Storage)
}

func TestConfigureOcspContractCertValidator(t *testing.T) {
	cfg := &config.DefaultConfig
	cfg.ContractCertValidator.Ocsp.RootCertProvider.File.FileNames = []string{"testdata/root_ca.pem"}

	settings, err := config.Configure(context.TODO(), cfg)
	require.NoError(t, err)
	require.NotNil(t, settings.ContractCertValidationService)
}

func TestConfigureOpcpContractCertProvider(t *testing.T) {
	_ = os.Setenv("TEST_OPCP_TOKEN", "test-token")
	defer func() {
		_ = os.Unsetenv("TEST_OPCP_TOKEN")
	}()

	cfg := &config.DefaultConfig
	cfg.ContractCertValidator.Ocsp.RootCertProvider.File.FileNames = []string{"testdata/root_ca.pem"}
	cfg.ContractCertProvider.Type = "opcp"
	cfg.ContractCertProvider.Opcp = &config.OpcpContractCertProviderConfig{
		Url: "http://localhost:8080",
		HttpAuth: config.HttpAuthConfig{
			Type: "env_token",
			EnvToken: &config.EnvHttpTokenConfig{
				EnvVar: "TEST_OPCP_TOKEN",
			},
		},
	}

	settings, err := config.Configure(context.TODO(), cfg)
	require.NoError(t, err)
	require.NotNil(t, settings.ContractCertProviderService)
}

func TestConfigureDefaultContractCertProvider(t *testing.T) {
	cfg := &config.DefaultConfig
	cfg.ContractCertValidator.Ocsp.RootCertProvider.File.FileNames = []string{"testdata/root_ca.pem"}
	cfg.ContractCertProvider.Type = "default"

	settings, err := config.Configure(context.TODO(), cfg)
	require.NoError(t, err)
	require.NotNil(t, settings.ContractCertProviderService)
}

func TestConfigureOpcpChargeStationCertProvider(t *testing.T) {
	_ = os.Setenv("TEST_OPCP_TOKEN", "test-token")
	defer func() {
		_ = os.Unsetenv("TEST_OPCP_TOKEN")
	}()

	cfg := &config.DefaultConfig
	cfg.ContractCertValidator.Ocsp.RootCertProvider.File.FileNames = []string{"testdata/root_ca.pem"}
	cfg.ChargeStationCertProvider.Type = "opcp"
	cfg.ChargeStationCertProvider.Opcp = &config.OpcpChargeStationCertProviderConfig{
		Url: "http://localhost:8080",
		HttpAuth: config.HttpAuthConfig{
			Type: "env_token",
			EnvToken: &config.EnvHttpTokenConfig{
				EnvVar: "TEST_OPCP_TOKEN",
			},
		},
	}

	settings, err := config.Configure(context.TODO(), cfg)
	require.NoError(t, err)
	require.NotNil(t, settings.ContractCertProviderService)
}

func TestConfigureDefaultChargeStationCertProvider(t *testing.T) {
	cfg := &config.DefaultConfig
	cfg.ContractCertValidator.Ocsp.RootCertProvider.File.FileNames = []string{"testdata/root_ca.pem"}
	cfg.ChargeStationCertProvider.Type = "default"

	settings, err := config.Configure(context.TODO(), cfg)
	require.NoError(t, err)
	require.NotNil(t, settings.ContractCertProviderService)
}

func TestConfigureKwHTariffService(t *testing.T) {
	cfg := &config.DefaultConfig
	cfg.ContractCertValidator.Ocsp.RootCertProvider.File.FileNames = []string{"testdata/root_ca.pem"}
	cfg.TariffService.Type = "kwh"

	settings, err := config.Configure(context.TODO(), cfg)
	require.NoError(t, err)
	require.NotNil(t, settings.ContractCertProviderService)
}
