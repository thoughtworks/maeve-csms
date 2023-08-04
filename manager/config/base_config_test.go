// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/manager/config"
	"testing"
)

func TestParseConfig(t *testing.T) {
	cfg := config.BaseConfig{}
	err := cfg.LoadFromFile("testdata/config.toml")
	require.NoError(t, err)

	want := config.BaseConfig{
		Api: config.ApiSettingsConfig{
			Addr: ":9410",
		},
		Mqtt: config.MqttSettingsConfig{
			Urls:   []string{"mqtt://127.0.0.1:1883"},
			Prefix: "cs",
			Group:  "manager",
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
						Cache: config.OpcpRootCertProviderCacheConfig{
							Ttl:  "24h",
							File: "/certs/root_certs.json",
						},
						HttpAuth: config.HttpAuthConfig{
							Type: "env_token",
							EnvToken: &config.EnvTokenHttpAuthConfig{
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
					EnvToken: &config.EnvTokenHttpAuthConfig{
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
					HubjectTestToken: &config.HubjectTestTokenHttpAuthConfig{
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
}
