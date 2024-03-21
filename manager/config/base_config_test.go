// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	clone "github.com/huandu/go-clone/generic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/config"
	"testing"
)

func TestParseConfig(t *testing.T) {
	cfg := clone.Clone(&config.DefaultConfig)
	err := cfg.LoadFromFile("testdata/config.toml")
	require.NoError(t, err)

	want := &config.BaseConfig{
		Api: config.ApiSettingsConfig{
			Addr:    ":9410",
			Host:    "example.com",
			WsPort:  80,
			WssPort: 443,
			OrgName: "Example",
		},
		Transport: config.TransportConfig{
			Type: "mqtt",
			Mqtt: &config.MqttSettingsConfig{
				Urls:              []string{"mqtt://127.0.0.1:1883"},
				Prefix:            "cs",
				Group:             "manager",
				ConnectTimeout:    "10s",
				ConnectRetryDelay: "1s",
				KeepAliveInterval: "10s",
			},
		},
		Ocpp: config.OcppSettingsConfig{
			HeartbeatInterval: "10m",
			Ocpp16Enabled:     false,
			Ocpp201Enabled:    true,
		},
		Observability: config.ObservabilitySettingsConfig{
			LogFormat:         "text",
			OtelCollectorAddr: "localhost:4317",
			TlsKeylogFile:     "/keylog/manager.log",
		},
		Storage: config.StorageConfig{
			Type: "firestore",
			FirestoreStorage: &config.FirestoreStorageConfig{
				ProjectId: "*detect-project-id*",
			},
		},
		ContractCertValidator: config.ContractCertValidatorConfig{
			Type: "ocsp",
			Ocsp: &config.OcspContractCertValidatorConfig{
				RootCertProvider: config.RootCertProviderConfig{
					Type: "opcp",
					Opcp: &config.OpcpRootCertProviderConfig{
						Url: "https://open.plugncharge-test.hubject.com",
						Ttl: "24h",
						HttpAuth: config.HttpAuthConfig{
							Type: "env_token",
							EnvToken: &config.EnvHttpTokenConfig{
								EnvVar: "RCP_TOKEN",
							},
						},
					},
				},
				MaxAttempts: 3,
			},
		},
		ContractCertProvider: config.ContractCertProviderConfig{
			Type: "opcp",
			Opcp: &config.OpcpContractCertProviderConfig{
				Url: "https://open.plugncharge-test.hubject.com",
				HttpAuth: config.HttpAuthConfig{
					Type: "env_token",
					EnvToken: &config.EnvHttpTokenConfig{
						EnvVar: "EST_TOKEN",
					},
				},
			},
		},
		ChargeStationCertProvider: config.ChargeStationCertProviderConfig{
			Type: "opcp",
			Opcp: &config.OpcpChargeStationCertProviderConfig{
				Url: "https://open.plugncharge-test.hubject.com",
				HttpAuth: config.HttpAuthConfig{
					Type: "hubject_test_token",
					HubjectTestToken: &config.HubjectTestHttpTokenConfig{
						Url: "https://hubject.stoplight.io/docs/open-plugncharge/6bb8b3bc79c2e-authorization-token",
					},
				},
			},
		},
		TariffService: config.TariffServiceConfig{
			Type: "kwh",
		},
	}

	assert.Equal(t, want, cfg)

	err = cfg.Validate()
	assert.NoError(t, err)
}

func TestValidateConfig(t *testing.T) {
	cfg := clone.Clone(&config.DefaultConfig)
	err := cfg.Validate()
	assert.NoError(t, err)
}
